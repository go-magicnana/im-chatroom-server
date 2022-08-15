package handler

import (
	"im-chatroom-broker/ctx"
	err "im-chatroom-broker/error"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/service"
	"im-chatroom-broker/util"
	"im-chatroom-broker/zaplog"
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

func (s SignalHandler) Handle(c *ctx.Context, packet *protocol.Packet) (*protocol.Packet, error) {

	ret := protocol.NewResponseError(packet, err.TypeNotAllow)

	switch packet.Header.Type {
	case protocol.TypeSignalPing:
		return ping(c, packet)

	case protocol.TypeSignalLogin:
		a := protocol.JsonSignalLogin(packet.Body)
		return login(c, packet, a)

	case protocol.TypeSignalJoinRoom:
		a := protocol.JsonSignalJoinRoom(packet.Body)
		return joinRoom(c, packet, a)

	case protocol.TypeSignalLeaveRoom:
		a := protocol.JsonSignalLeaveRoom(packet.Body)
		return leaveRoom(c, packet, a)

	case protocol.TypeSignalChangeRoom:
		a := protocol.JsonSignalChangeRoom(packet.Body)
		return changeRoom(c, packet, a)
	}
	return ret, nil
}

func ping(c *ctx.Context, packet *protocol.Packet) (*protocol.Packet, error) {
	service.RefreshUserClient(c.UserId)
	service.RefreshUserInfo(c.UserId)
	return nil, nil
}

func login(c *ctx.Context, packet *protocol.Packet, body *protocol.MessageBodySignalLogin) (*protocol.Packet, error) {

	zaplog.Logger.Debugf("Handler login %s %s", c.ClientName, body.Token)

	token := body.Token
	//device := body.Device

	user, e := service.GetUserAuth(token)

	if e != nil {
		return protocol.NewResponseError(packet, err.Unauthorized), nil
	}

	//exist, _ := service.GetUserDevice(ctx, c.ClientName())

	//if exist != nil && util.IsNotEmpty(exist.State) {
	//	if exist.State == strconv.FormatInt(int64(context2.Login), 10) {
	//
	//		alreadyLogin(ctx, c.ClientName(), packet)
	//	}
	//}

	c.UserId = user.UserId
	//service.SetUserClient(ctx, user.UserId, conn.RemoteAddr().String())

	//devices := service.GetUserClients(ctx, user.UserId)
	//
	userInfo := protocol.UserInfo{
		UserId: user.UserId,
		Token:  token,
		Name:   user.Name,
		Avatar: user.Avatar,
		Gender: user.Gender,
		Role:   user.Role,
	}
	service.SetUserInfo(userInfo)

	service.SetUserClient(c.UserId, c.ClientName)

	//userDevice := protocol.UserDevice{
	//	ClientName: c.ClientName(),
	//	UserId:     user.UserId,
	//	Device:     device,
	//	Broker:     c.Broker(),
	//}
	//service.SetUserDevice(ctx, userDevice)
	//
	//service.SetUserDevice2Login(ctx, c.ClientName(), context2.Login)

	//service.SetUserContext(c.ClientName(), c)

	//service.SetBrokerClients(ctx, conn.LocalAddr().String(), conn.RemoteAddr().String())

	loginUser := protocol.MessageBodySignalLoginRes{
		User: userInfo,
	}
	p := protocol.NewResponseOK(packet, loginUser)

	service.DelUserAuth(token)

	zaplog.Logger.Debugf("Handler login %s %s %s", c.ClientName, body.Token, user.UserId)

	return p, nil
}

func joinRoom(c *ctx.Context, packet *protocol.Packet, body *protocol.MessageBodySignalJoinRoom) (*protocol.Packet, error) {

	zaplog.Logger.Debugf("Handler join %s %s %s", c.ClientName, body.UserId, body.RoomId)

	if util.IsEmpty(body.RoomId) {
		return protocol.NewResponseError(packet, err.InvalidRequest.Format("roomId")), nil
	}

	//service.SetUserDevice2InRoom(ctx, c.ClientName(), body.RoomId)

	//service.SetRoomUser(ctx, body.RoomId, c.ClientName())

	c.RoomId = body.RoomId
	c.UserId = body.UserId

	service.SetRoomClients(c.Broker, c.RoomId, c.ClientName)

	noticeJoinRoom(packet.Header.MessageId, body.UserId, body.RoomId)

	//body.RoomBlocked = service.GetRoomBlocked(ctx, body.RoomId)
	//body.Blocked = service.GetRoomMemberBlocked(ctx, body.RoomId, body.UserId)

	zaplog.Logger.Debugf("Handler join %s %s %s", c.ClientName, body.UserId, body.RoomId)

	return protocol.NewResponseOK(packet, body), nil
}

func leaveRoom(c *ctx.Context, packet *protocol.Packet, a *protocol.MessageBodySignalLeaveRoom) (*protocol.Packet, error) {

	zaplog.Logger.Debugf("Handler leave %s %s %s", c.ClientName, a.UserId, a.RoomId)

	service.RemRoomClients(c.Broker, c.RoomId, c.ClientName)

	c.RoomId = ""
	c.UserId = a.UserId

	//userDevice, _ := service.GetUserDevice(ctx, c.ClientName())
	//service.DelUserDeviceInRoom(ctx, c.ClientName())
	//
	//service.DelRoomUser(ctx, userDevice.RoomId, c.ClientName())

	noticeLeaveRoom(packet.Header.MessageId, a.UserId, a.RoomId)

	zaplog.Logger.Debugf("Handler leave %s %s %s", c.ClientName, a.UserId, a.RoomId)

	return protocol.NewResponseOK(packet, nil), nil
}

func changeRoom(c *ctx.Context, packet *protocol.Packet, body *protocol.MessageBodySignalChangeRoom) (*protocol.Packet, error) {

	zaplog.Logger.Debugf("Handler leave %s %s %s %s", c.ClientName, body.UserId, body.OldRoomId, body.NewRoomId)

	if util.IsEmpty(body.NewRoomId) {
		return protocol.NewResponseError(packet, err.InvalidRequest.Format("NewRoomId")), nil
	}

	service.RemRoomClients(c.Broker, c.RoomId, c.ClientName)

	c.RoomId = body.NewRoomId
	c.UserId = body.UserId

	service.SetRoomClients(c.Broker, c.RoomId, c.ClientName)

	//info, _ := service.GetUserDevice(ctx, c.ClientName())
	//service.DelRoomUser(ctx, info.RoomId, c.ClientName())
	//service.SetRoomUser(ctx, body.RoomId, c.ClientName())
	//
	//service.SetUserDevice2InRoom(ctx, c.ClientName(), body.RoomId)

	noticeLeaveRoom(packet.Header.MessageId, body.UserId, body.OldRoomId)
	noticeJoinRoom(packet.Header.MessageId, body.UserId, body.NewRoomId)

	//body.RoomBlocked = service.GetRoomBlocked(ctx, body.NewRoomId)
	//body.Blocked = service.GetRoomMemberBlocked(ctx, body.OldRoomId, body.NewRoomId)

	zaplog.Logger.Debugf("Handler leave %s %s %s %s", c.ClientName, body.UserId, body.OldRoomId, body.NewRoomId)

	return protocol.NewResponseOK(packet, body), nil
}
