package handler

import (
	"golang.org/x/net/context"
	err "im-chatroom-broker/error"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/thread"
	"net"
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

func (d CustomHandler) Handle(ctx context.Context, conn net.Conn, packet *protocol.Packet, c *thread.ConnectClient) (*protocol.Packet, error) {
	ret := protocol.NewResponseError(packet, err.TypeNotAllow)

	switch packet.Header.Type {
	case protocol.TypeCustomNone:
		return custom(ctx, conn, packet, c)
	}
	return ret, nil
}

func custom(ctx context.Context, conn net.Conn, packet *protocol.Packet, c *thread.ConnectClient) (*protocol.Packet, error) {

	return deliver(ctx, conn, packet, c)
}
