package client

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/redis"
	"im-chatroom-broker/serializer"
	"im-chatroom-broker/service"
	"im-chatroom-broker/util"
	"im-chatroom-broker/zaplog"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup

func Start(role, serverIp string) {

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
	time.Sleep(time.Second * 10)

	//
	//
	//go sendPing(conn)

	if "send" == role {
		sendMsg(conn)
	}
	wg.Wait()
}

func read(conn net.Conn) {

	//filePath := "~/work/" + conn.LocalAddr().String() + ".txt"
	//file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	//if err != nil {
	//	fmt.Println("文件打开失败", err)
	//}
	//defer file.Close()

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

		p, e := serializer.DecodePacket(body)

		if e != nil || p == nil {
			return
		}

		zaplog.Logger.Debugf("ReadOK %s %s C:%d T:%d F:%d %s", conn.RemoteAddr().String(), p.Header.MessageId, p.Header.Command, p.Header.Type, p.Header.Flow, p.Body)

		if p.Header.Command == protocol.CommandContent && p.Header.Flow == protocol.FlowDeliver {
			text := protocol.JsonContentText(p.Body)
			service.AddUserClientMessage(context.Background(),conn.LocalAddr().String(),text.Content)
		}

	}

}

func write(conn net.Conn, p *protocol.Packet) error {

	j := serializer.SingleJsonSerializer()

	bs, e := j.EncodePacket(p)
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

	//SetUserAuth(i)

	header := protocol.MessageHeader{
		MessageId: "LoginMessageId-" + randCreator(8),
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
		MessageId: "JoinRoomMessageId-" + randCreator(8),
		Command:   protocol.CommandSignal,
		Flow:      protocol.FlowUp,
		Type:      protocol.TypeSignalJoinRoom,
	}

	body := protocol.MessageBodySignalJoinRoom{
		RoomId: "1",
		UserId: "1001",
	}

	packet := protocol.Packet{
		Header: header, Body: body,
	}

	write(conn, &packet)

}

func randCreator(l int) string {
	str := "0123456789abcdefghigklmnopqrstuvwxyz"
	strList := []byte(str)

	result := []byte{}
	i := 0

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i < l {
		new := strList[r.Intn(len(strList))]
		result = append(result, new)
		i = i + 1
	}
	return string(result)
}

func sendPing(conn net.Conn) {

	header := protocol.MessageHeader{
		MessageId: "PingMessageId-" + randCreator(8),
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
	for i := 0; i < 100; i++ {

		header := protocol.MessageHeader{
			MessageId: "ContentMessageId-" + randCreator(8),
			Command:   protocol.CommandContent,
			Flow:      protocol.FlowUp,
			Type:      protocol.TypeContentText,
			Target:    protocol.TargetRoom,
			To:        "1",
		}

		body := protocol.MessageBodyContentText{
			Content: strconv.Itoa(i) + "\n",
		}

		packet := protocol.Packet{
			Header: header, Body: body,
		}

		write(conn, &packet)
	}

}

func writeFile(path string) *os.File {
	fi, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	if err != nil {
		util.Panic(err)
	}
	defer fi.Close()

	//创建新Writer，其缓冲区有默认大小
	//writer := bufio.NewWriter(fi)
	//将信息写入缓存
	//_,err = writer.WriteString(info)
	//if err != nil {
	//	return
	//}
	////将缓存数据写入文件
	//err = writer.Flush()
	//if err != nil {
	//	return
	//}
	return fi
}

func SetUserAuth(userId string) error {

	u := protocol.UserInfo{

		UserId: userId,
		Token:  userId,
		Name:   "name-" + userId,
	}

	data, _ := json.Marshal(u)

	redis.Rdb.Set(context.Background(), service.UserAuth+u.Token, data, time.Minute*30)

	return nil

}
