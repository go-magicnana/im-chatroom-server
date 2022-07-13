package serializer

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"im-chatroom-broker/protocol"
	"io"
	"testing"
)

func TestJsonSerializer_Decode(t *testing.T) {

	header := protocol.MessageHeader{
		MessageId: "e10adc3949ba59abbe56e057f20f883e",
		Command:   protocol.CommandSignal,
		Flow:      protocol.FlowUp,
		Type:      protocol.TypeSignalPing,
	}

	body := protocol.MessageBodySignalPing{
		UserId: "1001",
	}

	packet := protocol.Packet{
		header, body,
	}


	msg, _ := json.Marshal(packet)
	l1 := len(msg)
	l2 := uint32(l1)


	buffer := new(bytes.Buffer) //直接使用 new 初始化，可以直接使用
	binary.Write(buffer, binary.BigEndian, Version)
	binary.Write(buffer, binary.BigEndian, l2)
	buffer.Write(msg)

	fmt.Println(len(buffer.Bytes()))



	var version uint8
	binary.Read(buffer, binary.BigEndian, &version)

	var length uint32
	binary.Read(buffer, binary.BigEndian, &length)


	b := make([]uint8, length)
	io.ReadFull(buffer, b)

	fmt.Println("ss")


}

func TestJsonSerializer_Encode(t *testing.T) {

}
