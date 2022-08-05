package serializer

import (
	"errors"
	"im-chatroom-broker/protocol"
	"net"
)

var ErrIncompletePacket = errors.New("incomplete packet")
var ErrInvalidPacket = errors.New("invalid packet")

const (
	Version = byte(1)
)

type Serializer interface {
	Name() (string, error)
	EncodePacket(packet *protocol.Packet) ([]byte, error)
	DecodePacket(bytes []byte) (*protocol.Packet, error)
	Version() byte
	Write(conn net.Conn, p *protocol.Packet) error
}
