package handler

import (
	"golang.org/x/net/context"
	context2 "im-chatroom-broker/context"
	"im-chatroom-broker/protocol"
)

type Handler interface {
	Handle(ctx context.Context, c *context2.Context, packet *protocol.Packet) *protocol.Packet
}
