package handler

import (
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/server"
)

type Handler interface {
	Handle(packet *protocol.Packet,c *server.Context)
}