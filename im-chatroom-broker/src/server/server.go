package server

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	context2 "im-chatroom-broker/context"
	"im-chatroom-broker/handler"
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

	//zaplog.InitLogger()
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

	ip, e := util.ExternalIP()
	if e != nil {
		util.Panic(e)
	}

	brokerAddress := ip.String() + addr

	//SetBroker(ctx,1, address)

	//UserLocal2String()

	for {
		select {
		case <-ctx.Done():
			return
		default:

			fmt.Println(util.CurrentSecond(), "Accept 等待客户端连接")
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

			//setDirtyConnection(c)

			go read(ctx, cancel, c)

			//go printConn()

			//go Read(c)
		}
	}
}

func catch(e error) {

}

func read(ctx context.Context, cancel context.CancelFunc, c *context2.Context) {

	defer c.Conn.Close()

	//defer c.CancelFunc()
	//SetReadDeadlineOnCancel(c.Ctx, c.CancelFunc, c.Conn)

	serializer := serializer.SingleJsonSerializer()

	for {

		fmt.Println(util.CurrentSecond(), "Read 等待客户端写入")

		c.Conn.SetReadDeadline(time.Now().Add(time.Second * 10))

		meta := make([]byte, protocol.MetaVersionBytes+protocol.MetaLengthBytes)
		ml, me := c.Conn.Read(meta)

		fmt.Println(me)

		switch me.(type) {
		case *net.OpError:
			operror := me.(*net.OpError)
			fmt.Println("operror",operror)
		}

		if me == io.EOF {
			//c.CancelFunc()
			break
		}

		if me != nil {
			continue
		}

		if ml != protocol.MetaVersionBytes+protocol.MetaLengthBytes {
			continue
		}

		version := meta[0]

		if version != serializer.Version() {
			continue
		}

		length := binary.BigEndian.Uint32(meta[1:])
		body := make([]byte, length)
		c.Conn.Read(body)

		packet, e := serializer.DecodePacket(body, c)

		if e != nil || packet == nil {
			return
		}

		fmt.Println(util.CurrentSecond(), "Read 读取客户端写入", packet)

		go process(ctx, cancel, c, packet)
	}

	fmt.Println("read thread exit")

}

func write(p *protocol.Packet, c *context2.Context) error {
	//SetWriteDeadlineOnCancel(c.Ctx, c.CancelFunc, c.Conn)

	serializer := serializer.SingleJsonSerializer()

	bs, e := serializer.EncodePacket(p, c)
	if bs == nil {
		return errors.New("empty packet")
	}

	if e != nil {
		return e
	}

	buffer := new(bytes.Buffer)

	binary.Write(buffer, binary.BigEndian, serializer.Version())

	length := uint32(len(bs))
	binary.Write(buffer, binary.BigEndian, length)

	buffer.Write(bs)
	_, err := c.Conn.Write(buffer.Bytes())

	fmt.Println(util.CurrentSecond(), "Write 等待客户端读取", p)

	if err != nil {
		return errors.New("write response error +" + err.Error())
	} else {
		return nil
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
	switch packet.Header.Command {
	case protocol.CommandDefault:
		ret = handler.SingleDefaultHandler().Handle(ctx, c, packet)
		break
	case protocol.CommandSignal:
		ret = handler.SingleSignalHandler().Handle(ctx, c, packet)
	}

	write(ret, c)

}

func goHandle(s string, context *context2.Context) {

	//fmt.Println(strings.Contains(s,"\n"))
	//
	//if startWith(s,"auth"){
	//	userId := s[5:]
	//
	//	v,e:=userMap.Load(userId)
	//	if e{
	//		context.UserId = v.(string)
	//		doAuth(context.UserId,context)
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
