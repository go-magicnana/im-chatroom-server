package handler

import (
	"im-chatroom-broker/context"
	err "im-chatroom-broker/error"
	"im-chatroom-broker/protocol"
	"sync"
)

var onceDefaultHandler sync.Once

var defaultHandler *DefaultHandler

func SingleDefaultHandler() *DefaultHandler {
	onceDefaultHandler.Do(func() {
		defaultHandler = &DefaultHandler{}
	})

	return defaultHandler
}

type DefaultHandler struct{}

func (d DefaultHandler) Handle(packet *protocol.Packet, c *context.Context) *protocol.Packet {
	ret := protocol.NewResponseError(packet, err.CommandNotAllow)
	return ret
}
