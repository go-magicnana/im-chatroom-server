package handler

import (
	"im-chatroom-broker/context"
	err "im-chatroom-broker/error"
	"im-chatroom-broker/protocol"
	"sync"
)

var onceSignalHandler sync.Once

var signalHandler *SignalHandler

func SingleSignalHandler() *SignalHandler {
	onceSignalHandler.Do(func() {
		signalHandler = &SignalHandler{}
	})

	return signalHandler
}

type SignalHandler struct{}

func (s SignalHandler) Handle(packet *protocol.Packet, c *context.Context) *protocol.Packet {

	ret := protocol.NewResponseError(packet, err.TypeNotAllow)

	/**
	TypeSignalPing       = 2101
	TypeSignalConnect    = 2102
	TypeSignalDisconnect = 2103
	TypeSignalJoinRoom   = 2104
	TypeSignalLeaveRoom  = 2105
	*/
	switch packet.Header.Type {
	case protocol.TypeSignalPing:

		body := protocol.JsonSignalPing(packet.Body, c)
		return ping(packet, body, c)
		break
	case protocol.TypeSignalConnect:

		a := protocol.JsonSignalConnect(packet.Body, c)
		return connect(packet, a, c)

		break
	case protocol.TypeSignalDisconnect:
		break
	case protocol.TypeSignalJoinRoom:
		a := protocol.JsonSignalJoinRoom(packet.Body, c)
		return joinRoom(packet, a, c)
	case protocol.TypeSignalLeaveRoom:
		break
	}

	return ret
}

func ping(packet *protocol.Packet, body *protocol.MessageBodySignalPing, c *context.Context) *protocol.Packet {
	user := protocol.User{
		UserId: body.UserId,
	}
	SetUser(user, c)
	return protocol.NewResponseOK(packet, nil)
}

func connect(packet *protocol.Packet, body *protocol.MessageBodySignalConnect, c *context.Context) *protocol.Packet {

	user := protocol.User{
		UserId: body.UserId,
		Name:   body.Name,
		Avatar: body.Avatar,
		Role:   body.Role,
	}

	SetUser(user, c)

	return protocol.NewResponseOK(packet, nil)
}

func joinRoom(packet *protocol.Packet, body *protocol.MessageBodySignalJoinRoom, c *context.Context) *protocol.Packet {

	return protocol.NewResponseOK(packet, nil)
}
