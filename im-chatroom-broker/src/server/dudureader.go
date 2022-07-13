package server

import (
	"encoding/binary"
	"fmt"
	"im-chatroom-broker/context"
	"im-chatroom-broker/util"
	"io"
	"net"
	"time"
)

func Read(c *context.Context) {
	_ = c.Conn.SetReadDeadline(time.Now().Add(time.Second * 30))

	//bs,_:=ioutil.ReadAll(c.Conn)
	//fmt.Println(len(bs))

	var version uint8
	binary.Read(c.Conn, binary.BigEndian, &version)

	var mask uint8
	binary.Read(c.Conn, binary.BigEndian, &mask)

	//var seq uint32
	//binary.Read(c.Conn, binary.BigEndian, &seq)

	var cmd uint16
	binary.Read(c.Conn, binary.BigEndian, &cmd)

	var length uint32
	binary.Read(c.Conn, binary.BigEndian, &length)

	body := make([]uint8, length)
	io.ReadFull(c.Conn, body)

	fmt.Println(version, mask, cmd, length, len(body))

	//userNew := pb.AuthMessage{}
	//err2 := proto.Unmarshal(body, &userNew)
	//if nil != err2 {
	//	fmt.Println(err2.Error())
	//}
	//
	//fmt.Println(userNew.MessageHead.From,userNew.Token)

	//length := readUint32(c.Conn)
	//fmt.Println(length)
	//
	//bs := readBody(length,c.Conn)
	//
	//userNew := pb.AuthMessage{}
	//err2 := proto.Unmarshal(bs, &userNew)
	//if nil != err2 {
	//	fmt.Println(err2.Error())
	//}
	//
	//fmt.Printf("%+v\n", userNew)

	//version := readUint8(c.Conn)
	//mark := readUint8(c.Conn)
	//seq := readUint8(c.Conn)
	//cmd := readUint16(c.Conn)
	//length := readUint32(c.Conn)//
	//
	//
	//fmt.Println(version, mark,seq, cmd, length)
	//bs := readBody(228,c.Conn)
	//
	//// 将二进制消息反序列化
	//userNew := pb.AuthMessage{}
	//err2 := proto.Unmarshal(bs, &userNew)
	//if nil != err2 {
	//	fmt.Println(err2.Error())
	//}
	//
	//fmt.Printf("%+v\n", userNew)

}

func readUint8(conn net.Conn, i *uint8) {
	binary.Read(conn, binary.LittleEndian, i)
}

func readUint16(conn net.Conn, i *uint16) {
	binary.Read(conn, binary.LittleEndian, i)
}

func readUint32(conn net.Conn, i *uint32) {
	binary.Read(conn, binary.LittleEndian, i)
}

//func readUint8(conn net.Conn) uint8 {
//	bs := make([]uint8, 1)
//	_, err := io.ReadFull(conn, bs)
//	if nil != err {
//		Panic(err)
//	}
//	return uint8(bs[0])
//}

//func readUint16(conn net.Conn) uint16 {
//	bs := make([]uint8, 2)
//	_, err := io.ReadFull(conn, bs)
//	if nil != err {
//		Panic(err)
//	}
//	//return uint16(bs[0]) | uint16(bs[1])<<8
//	return binary.LittleEndian.Uint16(bs)
//}
//
//func readUint32(conn net.Conn) uint32 {
//	bs := make([]uint8, 4)
//	_, err := io.ReadFull(conn, bs)
//	if nil != err {
//		Panic(err)
//	}
//	return binary.LittleEndian.Uint32(bs)
//}

func readBody(length uint32, conn net.Conn) []byte {
	bs := make([]uint8, length)
	_, err := io.ReadFull(conn, bs)
	if nil != err {
		util.Panic(err)
	}
	return bs
}
