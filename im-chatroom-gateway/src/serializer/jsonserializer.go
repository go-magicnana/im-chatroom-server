package serializer

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"im-chatroom-gateway/protocol"
	"im-chatroom-gateway/zaplog"
	"net"
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

func (j JsonSerializer) EncodePacket(packet *protocol.Packet) ([]byte, error) {
	bs, e := json.Marshal(packet)
	return bs, e

}

func (j JsonSerializer) DecodePacket(bytes []byte) (*protocol.Packet, error) {
	message := protocol.Packet{}
	e := json.Unmarshal(bytes, &message)
	return &message, e
}

func (j JsonSerializer) Write(conn net.Conn,p *protocol.Packet) error {

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

	zaplog.Logger.Debugf("WriteOK %s %s %d %d %s", conn.RemoteAddr().String(), p.Header.MessageId,p.Header.Command,p.Header.Type,p.Body)



	if err != nil {
		return errors.New("write response error +" + err.Error())
	} else {
		return nil
	}

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
