package handler

import (
	"encoding/json"
	"golang.org/x/net/context"
	context2 "im-chatroom-broker/context"
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

func (d CustomHandler) Handle(ctx context.Context, c *context2.Context, packet *protocol.Packet) (*protocol.Packet, error) {
	ret := protocol.NewResponseError(packet, err.TypeNotAllow)

	switch packet.Header.Type {
	case protocol.TypeCustomNone:
		return custom(ctx, c, packet)
	}
	return ret, nil
}

func custom(ctx context.Context, c *context2.Context, packet *protocol.Packet) (*protocol.Packet, error) {

	bytes, e := json.Marshal(packet.Body)
	if e != nil || len(bytes) == 0 {
		return nil, e
	}

	return deliver(ctx, c, packet)
}
