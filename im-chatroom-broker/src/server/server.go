package server

import (
	"context"
	"encoding/binary"
	"im-chatroom-broker/config"
	context2 "im-chatroom-broker/context"
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
	"sync"
	"time"
)

var counter = 100

//var wg sync.WaitGroup

var conns sync.Map

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

	brokerAddress := util.GetBrokerIp() + addr

	service.SetBrokerInstance(ctx, brokerAddress)
	service.SetBrokerAlive(ctx, brokerAddress)




	go service.AliveTask(ctx, brokerAddress)

	service.ProbeBroker(ctx)
	service.ProbeConn(ctx)
	service.ProbeRoom(ctx)

	//mq.Init()

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

			//if s.readDDL != 0 {
			//	_ = conn.SetReadDeadline(time.Now().Add(s.readDDL))
			//}
			//if s.writeDDL != 0 {
			//	_ = conn.SetWriteDeadline(time.Now().Add(s.writeDDL))
			//}

			ctx, cancel := context.WithCancel(ctx)

			c := context2.NewContext(brokerAddress, conn)

			c.Connect(conn.RemoteAddr().String())

			zaplog.Logger.Infof("Connected %s", conn.RemoteAddr())

			go read(ctx, cancel, c)

		}
	}
}


func read(ctx context.Context, cancel context.CancelFunc, c *context2.Context) {

	defer service.Close(ctx, c)

	serializer := serializer.SingleJsonSerializer()

	for {

		c.Conn().SetReadDeadline(time.Now().Add(time.Second * 60))

		meta := make([]byte, protocol.MetaVersionBytes+protocol.MetaLengthBytes)
		ml, me := c.Conn().Read(meta)

		switch me.(type) {
		case *net.OpError:
			if c.State() < context2.Login {

				zaplog.Logger.Errorf("ReadTimeOut %s Close Client", c.Conn().RemoteAddr())
				return
			} else {

				zaplog.Logger.Errorf("ReadTimeOut %s To Read Continue", c.Conn().RemoteAddr())
				continue
			}
		}

		if me == io.EOF {

			zaplog.Logger.Errorf("ReadClose %s Close Client", c.Conn().RemoteAddr())
			break
		}

		if me != nil {

			zaplog.Logger.Errorf("ReadError %s To Read Continue", c.Conn().RemoteAddr())
			continue
		}

		if ml != protocol.MetaVersionBytes+protocol.MetaLengthBytes {
			zaplog.Logger.Errorf("MetaError %s To Read Continue", c.Conn().RemoteAddr())
			continue
		}

		version := meta[0]

		if version != serializer.Version() {
			zaplog.Logger.Errorf("MetaOfVersionError %s To Read Continue", c.Conn().RemoteAddr())
			continue
		}

		length := binary.BigEndian.Uint32(meta[1:])
		body := make([]byte, length)
		c.Conn().Read(body)

		packet, e := serializer.DecodePacket(body, c)

		if e != nil || packet == nil {
			return
		}

		zaplog.Logger.Debugf("ReadOK %s Go Process %s %d %d %s", c.Conn().RemoteAddr(), packet.Header.MessageId, packet.Header.Command,packet.Header.Type,packet.Body)
		go process(ctx, cancel, c, packet)
	}

}

func process(ctx context.Context, cancel context.CancelFunc, c *context2.Context, packet *protocol.Packet) {

	var ret *protocol.Packet = nil
	var e error = nil
	switch packet.Header.Command {
	case protocol.CommandDefault:
		ret, e = handler.SingleDefaultHandler().Handle(ctx, c, packet)
		break
	case protocol.CommandSignal:
		ret, e = handler.SingleSignalHandler().Handle(ctx, c, packet)
		break
	case protocol.CommandContent:
		ret, e = handler.SingleContentHandler().Handle(ctx, c, packet)
		break
	case protocol.CommandCustom:
		ret, e = handler.CustomContentHandler().Handle(ctx, c, packet)
		break
	}

	if ret == nil {
		if e != nil {
			ret = protocol.NewResponseError(packet, err.Default)
		}
	}

	if ret != nil {
		serializer.SingleJsonSerializer().Write(c, ret)
	}

}
