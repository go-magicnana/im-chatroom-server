package handler

import (
	"im-chatroom-broker/ctx"
	err "im-chatroom-broker/error"
	"im-chatroom-broker/protocol"
	"sync"
)

var onceCustomHandler sync.Once

var customHandler *CustomHandler

func CustomContentHandler() *CustomHandler {
	onceCustomHandler.Do(func() {
		customHandler = &CustomHandler{}
	})

	return customHandler
}

type CustomHandler struct{}

func (d CustomHandler) Handle(c *ctx.Context, packet *protocol.Packet) (*protocol.Packet, error) {
	ret := protocol.NewResponseError(packet, err.TypeNotAllow)

	switch packet.Header.Type {
	case protocol.TypeCustomNone:
		return custom(c, packet)
	}
	return ret, nil
}

func custom(c *ctx.Context, packet *protocol.Packet) (*protocol.Packet, error) {

	return todeliver(c, packet)
}
