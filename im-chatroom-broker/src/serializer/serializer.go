package serializer

import (
	"im-chatroom-broker/context"
	"im-chatroom-broker/protocol"
)

const (
	Version = byte(1)
)

type Serializer interface {
	Name() (string, error)
	EncodePacket(packet *protocol.Packet, c *context.Context) ([]byte, error)
	DecodePacket(bytes []byte,c *context.Context) (*protocol.Packet, error)
	Version() byte
}
