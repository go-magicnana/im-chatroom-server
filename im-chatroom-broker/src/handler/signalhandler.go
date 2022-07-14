package handler

import (
	"golang.org/x/net/context"
	context2 "im-chatroom-broker/context"
	err "im-chatroom-broker/error"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/util"
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

	switch packet.Header.Type {
	case protocol.TypeSignalPing:
		return ping(ctx, c, packet)

	case protocol.TypeSignalLogin:
		a := protocol.JsonSignalLogin(packet.Body)
		return login(ctx, c, packet, a)

	case protocol.TypeSignalJoinRoom:
		a := protocol.JsonSignalJoinRoom(packet.Body)
		return joinRoom(ctx, c, packet, a)

	case protocol.TypeSignalLeaveRoom:
		return leaveRoom(ctx, c, packet)

	case protocol.TypeSignalChangeRoom:
		a := protocol.JsonSignalChangeRoom(packet.Body)
		return changeRoom(ctx, c, packet, a)
	}
	return ret, nil
}

func ping(ctx context.Context, c *context2.Context, packet *protocol.Packet) (*protocol.Packet, error) {
	c.Ping()
	return nil, nil
}

func login(ctx context.Context, c *context2.Context, packet *protocol.Packet, body *protocol.MessageBodySignalLogin) (*protocol.Packet, error) {

	if c.State() == context2.Login {
		return protocol.NewResponseError(packet, err.AlreadyLogin), nil
	}

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

	SetBrokerInfo(ctx, user.Broker, user.UserId)

	return protocol.NewResponseOK(packet, nil), nil
}

func joinRoom(ctx context.Context, c *context2.Context, packet *protocol.Packet, body *protocol.MessageBodySignalJoinRoom) (*protocol.Packet, error) {

	if util.IsEmpty(body.RoomId) {
		return protocol.NewResponseError(packet, err.InvalidRequest.Format("roomId")), nil
	}

	c.JoinRoom(body.RoomId)

	SetUserRoom(ctx, c.UserId(), body.RoomId)

	return protocol.NewResponseOK(packet, nil), nil
}

func leaveRoom(ctx context.Context, c *context2.Context, packet *protocol.Packet) (*protocol.Packet, error) {

	c.LeaveRoom()
	DelUserRoom(ctx, c.UserId())
	return protocol.NewResponseOK(packet, nil), nil
}

func changeRoom(ctx context.Context, c *context2.Context, packet *protocol.Packet, body *protocol.MessageBodySignalChangeRoom) (*protocol.Packet, error) {

	if util.IsEmpty(body.RoomId) {
		return protocol.NewResponseError(packet, err.InvalidRequest.Format("roomId")), nil
	}

	c.ChangeRoom(body.RoomId)
	SetUserRoom(ctx, c.UserId(), body.RoomId)

	return protocol.NewResponseOK(packet, nil), nil
}
