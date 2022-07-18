package server

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	context2 "im-chatroom-broker/context"
	err "im-chatroom-broker/error"
	"im-chatroom-broker/handler"
	"im-chatroom-broker/mq"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/serializer"
	"im-chatroom-broker/util"
	"io"
	"net"
	"strconv"
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

	handler.SetBrokerInstance(ctx, brokerAddress)

	go handler.BrokerAliveTask(ctx, brokerAddress)

	go consume(ctx)

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

			c.Connect()

			go read(ctx, cancel, c)

		}
	}
}

func consume(ctx context.Context) {
	mqConsumer := mq.NewRocketMqConsumer()
	mqConsumer.Subscribe("imchatroom_deliver", consumer.MessageSelector{},
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			for i := range msgs {
				fmt.Printf("subscribe callback : %v \n", msgs[i])
				go mqProcess(ctx, msgs[i])
			}
			return consumer.ConsumeSuccess, nil
		})

	mqConsumer.Start()
}

func mqProcess(ctx context.Context, ext *primitive.MessageExt) {
	tags := ext.GetTags()
	tagsInt, _ := strconv.Atoi(tags)
	packet := &protocol.Packet{}
	err := json.Unmarshal(ext.Body, packet)
	if err != nil {
		return
	}
	header := packet.Header
	if protocol.TargetRoom == tagsInt {
		room, _ := handler.GetRoom(ctx, header.To)
		for i := range room {
			userContext, exist := handler.GetUserContext(room[i])
			if !exist {
				continue
			}
			write(ctx, nil, userContext, packet)
		}
	} else {
		userContext, exist := handler.GetUserContext(header.To)
		if !exist {
			return
		}
		write(ctx, nil, userContext, packet)
	}
}

func close(ctx context.Context, cancel context.CancelFunc, c *context2.Context) {
	fmt.Println(util.CurrentSecond(), "Read 关闭线程 关闭连接")

	handler.DelUserInfo(ctx, c.UserKey())

	handler.DelUserContext(c.UserKey())

	handler.DelBrokerCapacity(ctx, c.Broker(), c.UserKey())

	c.Close()
}

func read(ctx context.Context, cancel context.CancelFunc, c *context2.Context) {

	defer close(ctx, cancel, c)

	serializer := serializer.SingleJsonSerializer()

	for {

		fmt.Println(util.CurrentSecond(), "Read 等待客户端写入")

		c.Conn().SetReadDeadline(time.Now().Add(time.Second * 9999))

		meta := make([]byte, protocol.MetaVersionBytes+protocol.MetaLengthBytes)
		ml, me := c.Conn().Read(meta)

		fmt.Println(me)

		switch me.(type) {
		case *net.OpError:
			if c.State() < context2.Login {
				return
			} else {
				continue
			}
		}

		if me == io.EOF {
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
		c.Conn().Read(body)

		packet, e := serializer.DecodePacket(body, c)

		if e != nil || packet == nil {
			return
		}

		fmt.Println(util.CurrentSecond(), "Read 读取客户端写入", packet.ToString())

		go process(ctx, cancel, c, packet)
	}

}

func write(ctx context.Context, cancel context.CancelFunc, c *context2.Context, p *protocol.Packet) error {
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
	_, err := c.Conn().Write(buffer.Bytes())

	fmt.Println(util.CurrentSecond(), "Write 等待客户端读取", p.ToString())

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

	if ret != nil {
		write(ctx, cancel, c, ret)
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
