package server

import (
	"context"
	"encoding/binary"
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

	addr := ":33121"

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

	go service.AliveTask(ctx, brokerAddress)

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

			c.Connect()

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

		zaplog.Logger.Debugf("ReadOK %s Go Process %s", c.Conn().RemoteAddr(), packet.ToString())
		go process(ctx, cancel, c, packet)
	}

}

//type ReadDeadliner interface {
//	SetReadDeadline(t time.Time) error
//}
//
//type WriteDeadliner interface {
//	SetWriteDeadline(t time.Time) error
//}
//
//func SetReadDeadlineOnCancel(ctx context.Context, cancel context.CancelFunc, d ReadDeadliner) {
//	go func() {
//		<-ctx.Done()
//		fmt.Println("receive done")
//		d.SetReadDeadline(time.Now())
//	}()
//}
//
//func SetWriteDeadlineOnCancel(ctx context.Context, cancel context.CancelFunc, d WriteDeadliner) {
//	go func() {
//		<-ctx.Done()
//		d.SetWriteDeadline(time.Now())
//	}()
//}

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
	}

	if ret == nil {
		if e != nil {
			ret = protocol.NewResponseError(packet, err.Default)
		}
	}

	//zaplog.Logger.Debugf("HandleOK %s %v", c.Conn().RemoteAddr(), ret.Header)

	if ret != nil {
		serializer.SingleJsonSerializer().Write(c, ret)
	}

}

func goHandle(s string, context *context2.Context) {

	//fmt.Println(strings.Contains(s,"\n"))
	//
	//if startWith(s,"auth"){
	//	userId := s[5:]
	//
	//	v,e:=userMap.Load(userId)
	//	if e{
	//		context.UserKey = v.(string)
	//		doAuth(context.UserKey,context)
	//	}
	//
	//}else{
	//
	//	if strings.Contains(s,":"){
	//		index := strings.Index(s,":")
	//		otherId := s[0:index]
	//		content := s[index:]+" \n"
	//
	//		v,e:=userMap.Load(otherId)
	//		if e{
	//			doDeliver(v.(string),content,context)
	//
	//		}
	//
	//	}
	//}
}

//func goConnection(context *Context) {
//
//	tmpBuffer := make([]byte, 0)
//
//	buffer := make([]byte, 1024)
//	messnager := make(chan byte)
//	for {
//		n, err := context.Conn.Read(buffer)
//		if err != nil {
//			Log(context.Conn.RemoteAddr().String(), " connection error: ", err)
//			return
//		}
//
//		tmpBuffer = protocol.Depack(append(tmpBuffer, buffer[:n]...))
//		Log("receive data string:", string(tmpBuffer))
//		TaskDeliver(tmpBuffer, context)
//		//start heartbeating
//		go HeartBeating(context, messnager, 10)
//		//check if get message from client
//		go GravelChannel(tmpBuffer, messnager)
//
//	}
//	defer context.Conn.Close()
//
//}
