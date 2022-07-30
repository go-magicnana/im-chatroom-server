package serializer

import (
	"im-chatroom-broker/protocol"
	"net"
)

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
