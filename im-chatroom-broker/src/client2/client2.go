package client2

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/serializer"
	"im-chatroom-broker/util"
	"im-chatroom-broker/zaplog"
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

	fmt.Println(util.CurrentSecond(), "Connect 连接服务端")

	go read(conn)

	time.Sleep(time.Second * 10)
	sendConnect(conn)

	//
	time.Sleep(time.Second * 10)
	sendJoinRoom(conn)
	//
	//
	time.Sleep(time.Second * 10)
	sendPing(conn)
	wg.Wait()
}

func read(conn net.Conn) {

	serializer := serializer.SingleJsonSerializer()

	for {

		fmt.Println(util.CurrentSecond(), "Read waiting server")

		meta := make([]byte, protocol.MetaVersionBytes+protocol.MetaLengthBytes)
		ml, me := conn.Read(meta)

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
		conn.Read(body)

		packet, e := serializer.DecodePacket(body, nil)

		if e != nil || packet == nil {
			return
		}

		fmt.Println(util.CurrentSecond(), "Read receive server", packet)

	}

}

func write(conn net.Conn, p *protocol.Packet) error {

	j := serializer.SingleJsonSerializer()

	bs, e := j.EncodePacket(p, nil)
	if bs == nil {
		return errors.New("empty packet")
	}

	if e != nil {
		return e
	}

	buffer := new(bytes.Buffer)

	binary.Write(buffer, binary.BigEndian, j.Version())

	length := uint32(len(bs))
	binary.Write(buffer, binary.BigEndian, length)

	buffer.Write(bs)
	_, err := conn.Write(buffer.Bytes())

	zaplog.Logger.Debugf("WriteOK %s %s %d %d %s", conn.RemoteAddr().String(), p.Header.MessageId, p.Header.Command, p.Header.Type, p.Body)

	if err != nil {
		return errors.New("write response error +" + err.Error())
	} else {
		return nil
	}

}

func sendConnect(conn net.Conn) {

	header := protocol.MessageHeader{
		MessageId: "e10adc3949ba59abbe56e057f20f883a",
		Command:   protocol.CommandSignal,
		Flow:      protocol.FlowUp,
		Type:      protocol.TypeSignalLogin,
	}

	body := protocol.MessageBodySignalLogin{
		Token:  "abcd",
		Device: "MAC",
		//UserId: "1001",
		//Name:   "张三丰",
		//Avatar: "https://img1.baidu.com/it/u=2848117662,2869906655&fm=253&fmt=auto&app=138&f=JPEG?w=500&h=501",
	}

	packet := protocol.Packet{
		Header: header, Body: body,
	}

	write(conn, &packet)

}

func sendJoinRoom(conn net.Conn) {

	header := protocol.MessageHeader{
		MessageId: "e10adc3949ba59abbe56e057f20f883b",
		Command:   protocol.CommandSignal,
		Flow:      protocol.FlowUp,
		Type:      protocol.TypeSignalJoinRoom,
	}

	body := protocol.MessageBodySignalJoinRoom{
		RoomId: "1",
	}

	packet := protocol.Packet{
		Header: header, Body: body,
	}

	write(conn, &packet)

}

func sendPing(conn net.Conn) {
	header := protocol.MessageHeader{
		MessageId: "e10adc3949ba59abbe56e057f20f883c",
		Command:   protocol.CommandSignal,
		Flow:      protocol.FlowUp,
		Type:      protocol.TypeSignalPing,
	}

	packet := protocol.Packet{
		Header: header, Body: nil,
	}

	write(conn, &packet)

	for {
		time.Sleep(time.Second * 10)
		write(conn, &packet)
	}

}
