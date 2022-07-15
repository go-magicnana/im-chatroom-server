package handler

import (
	"golang.org/x/net/context"
	context2 "im-chatroom-broker/context"
	err "im-chatroom-broker/error"
	"im-chatroom-broker/protocol"
	"sync"
)

var onceContentHandler sync.Once

var contentHandler *DefaultHandler

func SingleContentHandler() *DefaultHandler {
	onceDefaultHandler.Do(func() {
		defaultHandler = &DefaultHandler{}
	})

	return defaultHandler
}

type ContentHandler struct{}

/**
TypeContentText  = 4101
TypeContentEmoji = 4102
TypeContentReply = 4103
*/

func (d ContentHandler) Handle(ctx context.Context, c *context2.Context, packet *protocol.Packet) (*protocol.Packet, error) {
	ret := protocol.NewResponseError(packet, err.TypeNotAllow)

	switch packet.Header.Type {
	case protocol.TypeContentText:
		return ping(ctx, c, packet)

	case protocol.TypeContentEmoji:
		a := protocol.JsonSignalLogin(packet.Body)
		return login(ctx, c, packet, a)

	case protocol.TypeContentAt:
		a := protocol.JsonSignalJoinRoom(packet.Body)
		return joinRoom(ctx, c, packet, a)

	case protocol.TypeContentReply:
		return leaveRoom(ctx, c, packet)

	}
	return ret, nil
}
