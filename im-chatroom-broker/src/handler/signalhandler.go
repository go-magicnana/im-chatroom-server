package handler

import (
	"fmt"
	"golang.org/x/net/context"
	context2 "im-chatroom-broker/context"
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

func (s SignalHandler) Handle(ctx context.Context, c *context2.Context, packet *protocol.Packet) (*protocol.Packet, error) {

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
		return ping(ctx, c, packet, body),nil
		break
	case protocol.TypeSignalLogin:

		a := protocol.JsonSignalLogin(packet.Body, c)
		return connect(ctx, c, packet, a)

		break
	case protocol.TypeSignalJoinRoom:
		a := protocol.JsonSignalJoinRoom(packet.Body, c)
		return joinRoom(ctx, c, packet, a),nil
	case protocol.TypeSignalLeaveRoom:
		break
	}

	return ret,nil
}

func ping(ctx context.Context, c *context2.Context, packet *protocol.Packet, body *protocol.MessageBodySignalPing) *protocol.Packet {
	user := protocol.User{
		UserId: body.UserId,
	}

	fmt.Println(user)

	return protocol.NewResponseOK(packet, nil)
}

func connect(ctx context.Context, c *context2.Context, packet *protocol.Packet, body *protocol.MessageBodySignalLogin) (*protocol.Packet, error) {

	token := body.Token

	user, e := GetUserAuth(ctx, token)

	if e != nil {
		return protocol.NewResponseError(packet, err.Unauthorized), nil
	}

	user.Broker = c.Broker()
	user.Token = token
	c.Login(user.UserId, user.Token)

	SetUserInfo(ctx, user)

	SetUserContext(user, c)

	SetUserBroker(ctx, user.Broker, user.UserId)

	return protocol.NewResponseOK(packet, nil), nil
}

func joinRoom(ctx context.Context, c *context2.Context, packet *protocol.Packet, body *protocol.MessageBodySignalJoinRoom) *protocol.Packet {

	return protocol.NewResponseOK(packet, nil)
}
