package serializer

import (
	"encoding/json"
	"im-chatroom-broker/context"
	"im-chatroom-broker/protocol"
	"sync"
)

const Name = "JsonSerializer"

var once sync.Once

var jsonSerializer *JsonSerializer

func SingleJsonSerializer() *JsonSerializer {
	once.Do(func() {
		jsonSerializer = &JsonSerializer{}
	})

	return jsonSerializer
}

type JsonSerializer struct {
}

func (j JsonSerializer) Version() byte {
	return Version
}

func (j JsonSerializer) Name() (string, error) {
	return Name, nil
}

func (j JsonSerializer) EncodePacket(packet *protocol.Packet, c *context.Context) ([]byte, error) {
	bs, e := json.Marshal(packet)
	return bs, e

}

func (j JsonSerializer) DecodePacket(bytes []byte, c *context.Context) (*protocol.Packet, error) {

	//var version = make([]byte,1)
	//versionLength,versionE := c.Conn.Read(version)
	//
	//if versionLength != 1 || versionE!=nil{
	//	return nil,versionE
	//}
	//
	//
	//e:=binary.Read(c.Conn, binary.BigEndian, &version) //没有写入时，这里阻塞，写入后这里有值
	//if e!=nil {
	//	return nil,e
	//}
	//
	//
	//var length uint32
	//e = binary.Read(c.Conn, binary.BigEndian, &length)
	//if e!=nil {
	//	return nil,e
	//}
	//
	//body := make([]uint8, length)
	//size,e2:=io.ReadFull(c.Conn, body)
	//
	//fmt.Println(size)
	//if e2!=nil {
	//	return nil,e
	//}

	message := protocol.Packet{}

	e := json.Unmarshal(bytes, &message)

	return &message, e
}

//tmpBuffer := make([]byte, 0)
//
//buffer := make([]byte, 1024)
//messnager := make(chan byte)
//for {
//	n, err := context.Conn.Read(buffer)
//	if err != nil {
//		Log(context.Conn.RemoteAddr().String(), " connection error: ", err)
//		return
//	}
//
//	tmpBuffer = protocol.Depack(append(tmpBuffer, buffer[:n]...))
//	Log("receive data string:", string(tmpBuffer))
//	TaskDeliver(tmpBuffer, context)
//	//start heartbeating
//	go HeartBeating(context, messnager, 10)
//	//check if get message from client
//	go GravelChannel(tmpBuffer, messnager)
//
//}
//defer context.Conn.Close()
