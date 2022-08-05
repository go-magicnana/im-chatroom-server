package handler

import (
	"im-chatroom-broker/ctx"
	"im-chatroom-broker/protocol"
)

type Handler interface {
	Handle(c *ctx.Context, packet *protocol.Packet) (*protocol.Packet, error)
}
