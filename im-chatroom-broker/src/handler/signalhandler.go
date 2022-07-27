package handler

import (
	"golang.org/x/net/context"
	context2 "im-chatroom-broker/context"
	err "im-chatroom-broker/error"
	"im-chatroom-broker/mq"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/service"
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
	service.SetUserAlive(ctx, c.UserId(), c.ClientName())
	return nil, nil
}

func login(ctx context.Context, c *context2.Context, packet *protocol.Packet, body *protocol.MessageBodySignalLogin) (*protocol.Packet, error) {

	if c.State() == context2.Login {
		return protocol.NewResponseError(packet, err.AlreadyLogin), nil
	}

	token := body.Token
	device := body.Device

	user, e := service.GetUserAuth(ctx, token)

	if e != nil {
		return protocol.NewResponseError(packet, err.Unauthorized), nil
	}

	_, flag := c.Login(user.UserId)

	if !flag {
		return protocol.NewResponseError(packet, err.AlreadyLogin), nil
	}

	//exist, _ := service.GetUserDevice(ctx, c.ClientName())

	//if exist != nil && util.IsNotEmpty(exist.State) {
	//	if exist.State == strconv.FormatInt(int64(context2.Login), 10) {
	//
	//		alreadyLogin(ctx, c.ClientName(), packet)
	//	}
	//}

	service.SetUserClient(ctx, user.UserId, c.ClientName())

	devices := service.GetUserClients(ctx, user.UserId)

	userInfo := protocol.UserInfo{
		UserId: user.UserId,
		Token:  token,
		Device: devices,
		Name:   user.Name,
		Avatar: user.Avatar,
		Gender: user.Gender,
		Role:   user.Role,
	}
	service.SetUserInfo(ctx, userInfo)

	userDevice := protocol.UserDevice{
		ClientName: c.ClientName(),
		UserId:     user.UserId,
		Device:     device,
		Broker:     c.Broker(),
	}
	service.SetUserDevice(ctx, userDevice)

	service.SetUserDevice2Login(ctx, c.ClientName(), context2.Login)

	service.SetUserContext(&userDevice, c)

	service.SetBrokerCapacity(ctx, userDevice.Broker, c.ClientName())

	loginUser := protocol.MessageBodySignalLoginRes{
		User: userInfo,
	}
	p := protocol.NewResponseOK(packet, loginUser)

	service.DelUserAuth(ctx, token)

	return p, nil
}

func alreadyLogin(ctx context.Context, clientName string, packet *protocol.Packet) {
	p := protocol.Packet{
		Header: packet.Header,
		Body:   nil,
	}
	p.Header.Type = protocol.TypeSignalAlreadyLogin
	p.Header.Flow = protocol.FlowDeliver

	broker, _ := service.GetUserDeviceBroker(ctx, clientName)

	msg := protocol.PacketMessage{
		ClientName: clientName,
		Packet:     p,
	}
	mq.SendSync2One(broker, &msg)
}

func joinRoom(ctx context.Context, c *context2.Context, packet *protocol.Packet, body *protocol.MessageBodySignalJoinRoom) (*protocol.Packet, error) {

	if util.IsEmpty(body.RoomId) {
		return protocol.NewResponseError(packet, err.InvalidRequest.Format("roomId")), nil
	}

	c.JoinRoom(body.RoomId)

	service.SetUserDevice2InRoom(ctx, c.ClientName(), body.RoomId)

	service.SetRoomUser(ctx, body.RoomId, c.ClientName())

	noticeJoinRoom(ctx, c, body.RoomId)

	body.RoomBlocked = service.GetRoomBlocked(ctx, body.RoomId)
	body.Blocked = service.GetRoomMemberBlocked(ctx, body.RoomId, c.UserId())

	return protocol.NewResponseOK(packet, body), nil
}

func leaveRoom(ctx context.Context, c *context2.Context, packet *protocol.Packet) (*protocol.Packet, error) {

	c.LeaveRoom()
	userDevice, _ := service.GetUserDevice(ctx, c.ClientName())
	service.DelUserDeviceInRoom(ctx, c.ClientName())

	service.DelRoomUser(ctx, userDevice.RoomId, c.ClientName())

	noticeLeaveRoom(ctx, c, userDevice.RoomId)

	return protocol.NewResponseOK(packet, nil), nil
}

func changeRoom(ctx context.Context, c *context2.Context, packet *protocol.Packet, body *protocol.MessageBodySignalChangeRoom) (*protocol.Packet, error) {

	if util.IsEmpty(body.RoomId) {
		return protocol.NewResponseError(packet, err.InvalidRequest.Format("roomId")), nil
	}

	c.ChangeRoom(body.RoomId)

	info, _ := service.GetUserDevice(ctx, c.ClientName())
	service.DelRoomUser(ctx, info.RoomId, c.ClientName())
	service.SetRoomUser(ctx, body.RoomId, c.ClientName())

	service.SetUserDevice2InRoom(ctx, c.ClientName(), body.RoomId)

	noticeLeaveRoom(ctx, c, info.RoomId)
	noticeJoinRoom(ctx, c, body.RoomId)

	body.RoomBlocked = service.GetRoomBlocked(ctx, body.RoomId)
	body.Blocked = service.GetRoomMemberBlocked(ctx, body.RoomId, c.UserId())

	return protocol.NewResponseOK(packet, body), nil
}
