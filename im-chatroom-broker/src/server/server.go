package server

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/util"
	"io"
	"net"
	"sync"
	"time"
)

var counter = 100

var wg sync.WaitGroup

func Start() {

	userMap.Store("user1", "101")
	userMap.Store("user2", "102")
	userMap.Store("user3", "103")
	userMap.Store("user4", "104")
	userMap.Store("user5", "105")
	userMap.Store("user6", "106")
	userMap.Store("user7", "107")

	wg.Add(1)

	addr := ":33121"

	ctx := context.Background()

	go goListen(ctx, addr)

	wg.Wait()
}

func goListen(ctx context.Context, addr string) {

	netListen, _ := net.Listen("tcp", addr)
	defer netListen.Close()

	ip, e := util.ExternalIP()
	if e != nil {
		util.Panic(e)
	}

	address := ip.String() + addr

	//SetBroker(ctx,1, address)

	//UserLocal2String()

	for {
		select {
		case <-ctx.Done():
			return
		default:
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

			cctx, cancelFunc := context.WithCancel(ctx)

			c := NewContext(conn.RemoteAddr().String(), address, conn, cctx, cancelFunc)

			//setDirtyConnection(c)

			go goRead(c)

			//go Read(c)
		}
	}
}

func goRead(c *Context) {

	defer c.CancelFunc()
	SetReadDeadlineOnCancel(c.Ctx, c.Conn)

	messageIdBody := make([]uint8, 32)
	io.ReadFull(c.Conn, messageIdBody)

	var version uint8
	binary.Read(c.Conn, binary.BigEndian, &version)

	//var mask uint8
	//binary.Read(c.Conn, binary.BigEndian, &mask)

	//var seq uint32
	//binary.Read(c.Conn, binary.BigEndian, &seq)

	var cmd uint16
	binary.Read(c.Conn, binary.BigEndian, &cmd)

	var length uint32
	binary.Read(c.Conn, binary.BigEndian, &length)

	body := make([]uint8, length)
	io.ReadFull(c.Conn, body)

	message := protocol.Message{}

	json.Unmarshal(body, &message)

	packet := protocol.Packet{
		MessageId: util.B2s(messageIdBody),
		Version:   version,
		Command:   cmd,
		Message:   message,
	}

	go process(&packet, c)
}

func write(p *protocol.Packet, c *Context) {
	msg, _ := json.Marshal(p.Message)
	buffer := bytes.NewBufferString(p.MessageId)
	binary.Write(buffer, binary.LittleEndian, p.Version)
	binary.Write(buffer, binary.LittleEndian, p.Command)
	binary.Write(buffer, binary.LittleEndian, len(msg))
	bs := append(buffer.Bytes(), msg...)
	c.Conn.Write(bs)
}

type ReadDeadliner interface {
	SetReadDeadline(t time.Time) error
}

func SetReadDeadlineOnCancel(ctx context.Context, d ReadDeadliner) {
	go func() {
		<-ctx.Done()
		d.SetReadDeadline(time.Now())
	}()
}

func process(packet *protocol.Packet, c *Context) {

	switch packet.Command {
	case protocol.CommandDefault:
		p := protocol.NewResponseError(packet, &protocol.ErrorDefaultKey)
		write(p, c)
	}
}

func goHandle(s string, context *Context) {

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

var userMap sync.Map

func doAuth(userId string, context *Context) {
	SetUser(userId, context)
}

func doDeliver(userId string, content string, context *Context) {
	channel, exist := GetUserLocal(userId)
	if !exist {
		context.Conn.Write(util.S2b(userId + " offline"))
	} else {
		channel.Conn.Write(util.S2b(content))
	}
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
