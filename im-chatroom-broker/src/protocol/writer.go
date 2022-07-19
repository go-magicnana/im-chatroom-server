package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	context2 "im-chatroom-broker/context"
	"im-chatroom-broker/serializer"
	"im-chatroom-broker/util"
)

func Write(c *context2.Context, p *Packet) error {
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
