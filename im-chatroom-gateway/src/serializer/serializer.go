package serializer

import (
	"im-chatroom-gateway/src/context"
	"im-chatroom-gateway/src/protocol"
)

const (
	Version = byte(1)
)

type Serializer interface {
	Name() (string, error)
	EncodePacket(packet *protocol.Packet, c *context.Context) ([]byte, error)
	DecodePacket(bytes []byte, c *context.Context) (*protocol.Packet, error)
	Version() byte
	Write(c *context.Context, p *protocol.Packet) error
}
