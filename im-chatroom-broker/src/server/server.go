package server

import (
	"context"
	"encoding/binary"
	"errors"
	"im-chatroom-broker/config"
	"im-chatroom-broker/thread"
	"sync"

	//context2 "im-chatroom-broker/context"
	err "im-chatroom-broker/error"
	"im-chatroom-broker/handler"
	"im-chatroom-broker/service"
	"im-chatroom-broker/zaplog"
	//"im-chatroom-broker/mq"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/serializer"
	"im-chatroom-broker/util"
	"io"
	"net"
	"time"
)

var counter = 100

var channelMap sync.Map

func Start() {

	zaplog.InitLogger()

	//zaplog.Infof("Start ...")

	//wg.Add(1)

	addr := ":" + config.OP.Port

	ctx := context.Background()

	//go goListen(ctx, addr)
	listen(ctx, addr)

	//wg.Wait()
	//zaplog.Infof("Exit")

}

func listen(ctx context.Context, addr string) {

	netListen, err := net.Listen("tcp", addr)
	defer netListen.Close()

	util.Panic(err)

	var brokerAddress string
	if util.IsNotEmpty(config.OP.Ip) {
		brokerAddress = config.OP.Ip + addr
	} else {
		brokerAddress = util.GetBrokerIp() + addr
	}

	service.SetBrokerInstance(ctx, brokerAddress)
	//service.SetBrokerAlive(ctx, brokerAddress)

	go service.AliveTask(ctx, brokerAddress)
	//
	////service.ProbeBroker(ctx)
	//service.ProbeConn(ctx)
	//service.ProbeRoom(ctx)

	//mq.Init()

	limit := make(chan bool, 5000)

	for {
		select {
		case <-ctx.Done():
			return
		default:

			zaplog.Logger.Infof("Accept %s", brokerAddress)

			conn, err := netListen.Accept()
			if err != nil {
				util.Panic(err)
			}

			select {
			case limit <- true:

				zaplog.Logger.Infof("Connected %s", conn.RemoteAddr())

				ctx, cancel := context.WithCancel(ctx)

				channel := make(chan *protocol.InnerPacket, 65535)

				//c := context2.NewContext(brokerAddress, conn)

				//c.Connect(conn.RemoteAddr().String())

				cc := &thread.ConnectClient{
					Broker:     brokerAddress,
					ClientName: conn.RemoteAddr().String(),
					Channel:    channel,
					Conn:       conn,
				}

				thread.SetChannel(conn.RemoteAddr().String(), cc)

				go read(ctx, cancel, cc, conn)

				go write(ctx, cancel, cc, conn)
			default:
				zaplog.Logger.Infof("Overflow %s", conn.RemoteAddr())
				conn.Close()
			}
		}
	}
}

func readReturn(channel chan *protocol.InnerPacket, ip *protocol.InnerPacket) {

	if ip == nil {
		util.Panic(errors.New("nil response in read thread"))
	}

	if ip != nil {
		channel <- ip
	}
}

func read(
	ctx context.Context,
	cancel context.CancelFunc,
	cc *thread.ConnectClient,
	conn net.Conn) {

	serializer := serializer.SingleJsonSerializer()

	for {

		select {
		case <-ctx.Done():
			zaplog.Logger.Infof("ReadDone %s", conn.RemoteAddr().String())
			return
		default:
			//c.Readable()

			conn.SetReadDeadline(time.Now().Add(time.Second * 60))

			meta := make([]byte, protocol.MetaVersionBytes+protocol.MetaLengthBytes)
			ml, me := conn.Read(meta)

			//switch me.(type) {
			//case *net.OpError:
			//	if c.State() < context2.Login {
			//		zaplog.Logger.Errorf("ReadTimeOut %s Close Client", c.Conn().RemoteAddr())
			//		readReturn(c, protocol.NewQuit())
			//		return
			//	} else {
			//		zaplog.Logger.Errorf("ReadTimeOut %s To Read Continue", c.Conn().RemoteAddr())
			//		continue
			//	}
			//}

			if me == io.EOF {
				zaplog.Logger.Errorf("ReadClose %s Close Client", conn.RemoteAddr())
				readReturn(cc.Channel, protocol.NewQuit())
				return
			}

			if me != nil {
				zaplog.Logger.Errorf("ReadError %s To Read Continue", conn.RemoteAddr())
				continue
			}

			if ml != protocol.MetaVersionBytes+protocol.MetaLengthBytes {
				zaplog.Logger.Errorf("MetaError %s To Read Continue", conn.RemoteAddr())
				continue
			}

			version := meta[0]

			if version != serializer.Version() {
				zaplog.Logger.Errorf("MetaOfVersionError %s To Read Continue", conn.RemoteAddr())
				continue
			}

			length := binary.BigEndian.Uint32(meta[1:])
			body := make([]byte, length)
			conn.Read(body)

			packet, e := serializer.DecodePacket(body)

			if e != nil || packet == nil {
				readReturn(cc.Channel, protocol.NewQuit())
				return
			}

			//zaplog.Logger.Debugf("ReadOK %s %s C:%d T:%d F:%d %s", conn.RemoteAddr().String(), packet.Header.MessageId, packet.Header.Command, packet.Header.Type, packet.Header.Flow, packet.Body)

			go process(ctx, cancel, cc, packet, conn)

			//readReturn(c, protocol.NewResponse(res))
		}
	}

}

func close(ctx context.Context, cancel context.CancelFunc, conn net.Conn, cc *thread.ConnectClient) {

	service.RemUserClient(ctx, cc.UserId, conn.RemoteAddr().String())
	service.RemBrokerClients(ctx, conn.LocalAddr().String(), conn.RemoteAddr().String())
	thread.RemChannel(conn.RemoteAddr().String())
	thread.RemRoomChannel(cc.RoomId, conn.RemoteAddr().String())

	if cancel != nil {
		cancel()
	}

	conn.Close()

	zaplog.Logger.Infof("CloseByClient %s", conn.RemoteAddr().String())
}

func write(ctx context.Context,
	cancel context.CancelFunc,
	cc *thread.ConnectClient,
	conn net.Conn) {

	defer close(ctx, cancel, conn, cc)

	for {
		select {
		case <-ctx.Done():
			zaplog.Logger.Infof("WriteDone %s", conn.RemoteAddr())
			return
		case res := <-cc.Channel:

			//zaplog.Logger.Debugf("WriteCmd %s ", res.Cmd)

			if protocol.CmdQuit == res.Cmd {
				return
			} else {
				if res.Packet != nil {
					serializer.SingleJsonSerializer().Write(conn, res.Packet)
				}
			}
		default:
			continue
		}
	}

}

func process(ctx context.Context, cancel context.CancelFunc, c *thread.ConnectClient, packet *protocol.Packet, conn net.Conn) {

	var ret *protocol.Packet = nil
	var e error = nil
	switch packet.Header.Command {
	case protocol.CommandDefault:
		ret, e = handler.SingleDefaultHandler().Handle(ctx, conn, packet, c)
		break
	case protocol.CommandSignal:
		ret, e = handler.SingleSignalHandler().Handle(ctx, conn, packet, c)
		break
	case protocol.CommandContent:
		ret, e = handler.SingleContentHandler().Handle(ctx, conn, packet, c)
		break
	case protocol.CommandCustom:
		ret, e = handler.CustomContentHandler().Handle(ctx, conn, packet, c)
		break
	}

	if ret == nil {
		if e != nil {
			ret = protocol.NewResponseError(packet, err.Default)
		}
	}

	if ret != nil {
		c.Channel <- protocol.NewResponse(ret)
	}
}

//func (ch *ContextHolder) SetChannel(clientName string, channel chan *protocol.InnerPacket) {
//	ch.channelMap.Store(clientName, channel)
//}
//
//func (ch *ContextHolder) GetChannel(clientName string) chan *protocol.InnerPacket {
//	k, _ := ch.channelMap.Load(clientName)
//
//	return k.(chan *protocol.InnerPacket)
//
//}
//
//func (ch *ContextHolder) RemChannel(clientName string) {
//	ch.channelMap.Delete(clientName)
//}
//
//func (ch *ContextHolder) RanChannel(f func(key, value any) bool) {
//	ch.channelMap.Range(f)
//}
