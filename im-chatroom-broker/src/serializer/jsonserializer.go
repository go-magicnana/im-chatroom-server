package serializer

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/panjf2000/gnet/v2"
	"im-chatroom-broker/protocol"
	"net"
)

const Name = "JsonSerializer"

type JsonSerializer struct {
}

func NewJsonSerializer() *JsonSerializer {
	return &JsonSerializer{}
}

func (j *JsonSerializer) Version() byte {
	return Version
}

func (j *JsonSerializer) Name() (string, error) {
	return Name, nil
}

func (j JsonSerializer) EncodePacket(packet *protocol.Packet) ([]byte, error) {
	bs, e := json.Marshal(packet)
	return bs, e

}

func (j *JsonSerializer) DecodePacket(bytes []byte) (*protocol.Packet, error) {
	message := protocol.Packet{}
	e := json.Unmarshal(bytes, &message)
	return &message, e
}

func (j *JsonSerializer) Write(conn net.Conn, p *protocol.Packet) error {

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

	//zaplog.Logger.Debugf("WriteOK %s %s C:%d T:%d F:%d %s", conn.RemoteAddr().String(), p.Header.MessageId, p.Header.Command, p.Header.Type,p.Header.Flow, p.Body)

	if err != nil {
		return errors.New("write response error +" + err.Error())
	} else {
		return nil
	}

}

func (j *JsonSerializer) Encode(buf []byte) ([]byte, error) {

	buffer := new(bytes.Buffer)

	binary.Write(buffer, binary.BigEndian, j.Version())

	length := uint32(len(buf))
	binary.Write(buffer, binary.BigEndian, length)

	buffer.Write(buf)
	return buffer.Bytes(), nil
}

func (j *JsonSerializer) Decode(c gnet.Conn) ([]byte, error) {
	bodyOffset := 5
	buf, _ := c.Peek(bodyOffset)
	if len(buf) < bodyOffset {
		return nil, ErrIncompletePacket
	}

	if len(buf) != bodyOffset {
		return nil, ErrIncompletePacket

	}

	bodyLen := binary.BigEndian.Uint32(buf[1:bodyOffset])
	msgLen := bodyOffset + int(bodyLen)
	if c.InboundBuffered() < msgLen {
		return nil, ErrIncompletePacket
	}
	buf, _ = c.Peek(msgLen)
	_, _ = c.Discard(msgLen)

	return buf[bodyOffset:msgLen], nil
}
