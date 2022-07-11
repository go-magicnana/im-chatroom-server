package handler

import (
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/server"
	"im-chatroom-broker/util"
)

type CommandDefaultHandler struct{}

func (CommandDefaultHandler) Handle(packet *protocol.Packet, c *server.Context) {
	util.Panic("implement me")

}

