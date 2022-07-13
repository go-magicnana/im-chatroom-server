package handler

import (
	"im-chatroom-broker/context"
	"im-chatroom-broker/protocol"
)

type Handler interface {
	Handle(packet *protocol.Packet, c *context.Context) *protocol.Packet
}
