package client2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"im-chatroom-broker/context"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/serializer"
	"im-chatroom-broker/util"
	"net"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

func Start() {

	wg.Add(1)

	server := "localhost:33121"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	fmt.Println(util.CurrentSecond(),"Connect 连接服务端")

	go read(conn)

	time.Sleep(time.Second*10)
	sendConnect(conn)

	//
	time.Sleep(time.Second*10)
	sendJoinRoom(conn)
	//
	//
	time.Sleep(time.Second*10)
	sendPing(conn)
	wg.Wait()
}

func read(conn net.Conn) {

	c := context.Context{
		Conn: conn,
	}

	serializer := serializer.SingleJsonSerializer()

	for {


		fmt.Println(util.CurrentSecond(),"Read waiting server")


		meta := make([]byte, protocol.MetaVersionBytes+protocol.MetaLengthBytes)
		ml, me := c.Conn.Read(meta)

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

		packet, e := serializer.DecodePacket(body, &c)

		if e != nil || packet == nil {
			return
		}

		fmt.Println(util.CurrentSecond(),"Read receive server",packet)

	}

}

func write(packet *protocol.Packet, conn net.Conn) {

	c := context.Context{
		Conn: conn,
	}

	serializer := serializer.SingleJsonSerializer()

	bs, e := serializer.EncodePacket(packet, &c)

	if bs == nil || e != nil {
		return
	}

	buffer := new(bytes.Buffer)

	binary.Write(buffer, binary.BigEndian, serializer.Version())

	length := uint32(len(bs))
	binary.Write(buffer, binary.BigEndian, length)

	buffer.Write(bs)
	c.Conn.Write(buffer.Bytes())

	fmt.Println(util.CurrentSecond(),"Write send server",packet)


}

func sendConnect(conn net.Conn) {

	header := protocol.MessageHeader{
		MessageId: "e10adc3949ba59abbe56e057f20f883a",
		Command:   protocol.CommandSignal,
		Flow:      protocol.FlowUp,
		Type:      protocol.TypeSignalConnect,
	}

	body := protocol.MessageBodySignalConnect{
		UserId: "1001",
		Name:   "张三丰",
		Avatar: "https://img1.baidu.com/it/u=2848117662,2869906655&fm=253&fmt=auto&app=138&f=JPEG?w=500&h=501",
	}

	packet := protocol.Packet{
		Header: header, Body: body,
	}

	write(&packet, conn)

}

func sendJoinRoom(conn net.Conn) {

	header := protocol.MessageHeader{
		MessageId: "e10adc3949ba59abbe56e057f20f883b",
		Command:   protocol.CommandSignal,
		Flow:      protocol.FlowUp,
		Type:      protocol.TypeSignalJoinRoom,
	}

	body := protocol.MessageBodySignalJoinRoom{
		UserId: "1001",
		RoomId: "2001",
	}

	packet := protocol.Packet{
		Header: header, Body: body,
	}

	write(&packet, conn)

}

func sendPing(conn net.Conn) {
	header := protocol.MessageHeader{
		MessageId: "e10adc3949ba59abbe56e057f20f883c",
		Command:   protocol.CommandSignal,
		Flow:      protocol.FlowUp,
		Type:      protocol.TypeSignalPing,
	}

	body := protocol.MessageBodySignalPing{
		UserId: "1001",
	}

	packet := protocol.Packet{
		Header: header, Body: body,
	}

	write(&packet, conn)

	for{
		time.Sleep(time.Second*1)
		write(&packet, conn)
	}

}
