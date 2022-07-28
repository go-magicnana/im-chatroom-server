package client

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

func Start(serverIp string) {

	wg.Add(1)

	server := serverIp + ":33121"
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

	SendLogin(conn)
	time.Sleep(time.Second * 10)

	sendJoinRoom(conn)
	time.Sleep(time.Second * 2)

	//
	//
	go sendPing(conn)

	go sendMsg(conn)
	wg.Wait()
}

func read(conn net.Conn) {

	serializer := serializer.SingleJsonSerializer()

	for {

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

		p, e := serializer.DecodePacket(body, nil)

		if e != nil || p == nil {
			return
		}

		zaplog.Logger.Debugf("ReadOK %s %s C:%d T:%d F:%d %s", conn.RemoteAddr().String(), p.Header.MessageId, p.Header.Command, p.Header.Type, p.Header.Flow, p.Body)

		//if packet.Header.Command == protocol.CommandContent && packet.Header.Flow == protocol.FlowDeliver {
		//responseBody, _ := json.Marshal(packet.Body)
		//fmt.Println(util.CurrentSecond(), "Read receive server", string(responseBody))
		//}

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

	zaplog.Logger.Debugf("WriteOK %s %s C:%d T:%d F:%d %s", conn.RemoteAddr().String(), p.Header.MessageId, p.Header.Command, p.Header.Type, p.Header.Flow, p.Body)

	if err != nil {
		return errors.New("write response error +" + err.Error())
	} else {
		return nil
	}

}

func SendLogin(conn net.Conn) {

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

	for {
		write(conn, &packet)
		time.Sleep(time.Second * 10)
	}

}

func sendMsg(conn net.Conn) {
	header := protocol.MessageHeader{
		MessageId: "e10adc3949ba59abbe56e057f20f883D",
		Command:   protocol.CommandContent,
		Flow:      protocol.FlowUp,
		Type:      protocol.TypeContentText,
		Target:    protocol.TargetRoom,
	}

	body := protocol.MessageBodyContentText{
		Content: "你好 世界",
	}

	packet := protocol.Packet{
		Header: header, Body: body,
	}

	for {
		write(conn, &packet)
		time.Sleep(time.Second * 10)
	}

}
