package handler

import (
	"golang.org/x/net/context"
	err "im-chatroom-broker/error"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/service"
	"im-chatroom-broker/thread"
	"im-chatroom-broker/util"
	"net"
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

func (s SignalHandler) Handle(ctx context.Context, conn net.Conn, packet *protocol.Packet, c *thread.ConnectClient) (*protocol.Packet, error) {

	ret := protocol.NewResponseError(packet, err.TypeNotAllow)

	switch packet.Header.Type {
	case protocol.TypeSignalPing:
		//a := protocol.JsonSignalPing(packet.Body)
		return ping(ctx, conn, packet, c)

	case protocol.TypeSignalLogin:
		a := protocol.JsonSignalLogin(packet.Body)
		return login(ctx, conn, packet, a, c)

	case protocol.TypeSignalJoinRoom:
		a := protocol.JsonSignalJoinRoom(packet.Body)
		return joinRoom(ctx, conn, packet, a, c)

	case protocol.TypeSignalLeaveRoom:
		a := protocol.JsonSignalLeaveRoom(packet.Body)
		return leaveRoom(ctx, conn, packet, a, c)

	case protocol.TypeSignalChangeRoom:
		a := protocol.JsonSignalChangeRoom(packet.Body)
		return changeRoom(ctx, conn, packet, a, c)
	}
	return ret, nil
}

func ping(ctx context.Context, conn net.Conn, packet *protocol.Packet, c *thread.ConnectClient) (*protocol.Packet, error) {
	//service.RefreshUserClient(ctx,c.UserId())
	//service.RefreshUserInfo(ctx,c.UserId())
	return nil, nil
}

func login(ctx context.Context, conn net.Conn, packet *protocol.Packet, body *protocol.MessageBodySignalLogin, c *thread.ConnectClient) (*protocol.Packet, error) {

	token := body.Token
	//device := body.Device

	user, e := service.GetUserAuth(ctx, token)

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
	service.SetUserClient(ctx, user.UserId, conn.RemoteAddr().String())

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
	service.SetUserInfo(ctx, userInfo)

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

	service.SetBrokerClients(ctx, conn.LocalAddr().String(), conn.RemoteAddr().String())

	loginUser := protocol.MessageBodySignalLoginRes{
		User: userInfo,
	}
	p := protocol.NewResponseOK(packet, loginUser)

	service.DelUserAuth(ctx, token)

	return p, nil
}

func joinRoom(ctx context.Context, conn net.Conn, packet *protocol.Packet, body *protocol.MessageBodySignalJoinRoom, c *thread.ConnectClient) (*protocol.Packet, error) {

	if util.IsEmpty(body.RoomId) {
		return protocol.NewResponseError(packet, err.InvalidRequest.Format("roomId")), nil
	}

	//service.SetUserDevice2InRoom(ctx, c.ClientName(), body.RoomId)

	//service.SetRoomUser(ctx, body.RoomId, c.ClientName())

	c.RoomId = body.RoomId
	c.UserId = body.UserId
	thread.SetRoomChannels(body.RoomId, conn.RemoteAddr().String())

	//noticeJoinRoom(ctx, conn, packet.Header.MessageId, body.UserId, body.RoomId)

	//body.RoomBlocked = service.GetRoomBlocked(ctx, body.RoomId)
	//body.Blocked = service.GetRoomMemberBlocked(ctx, body.RoomId, body.UserId)

	return protocol.NewResponseOK(packet, body), nil
}

func leaveRoom(ctx context.Context, conn net.Conn, packet *protocol.Packet, a *protocol.MessageBodySignalLeaveRoom, c *thread.ConnectClient) (*protocol.Packet, error) {

	thread.RemRoomChannel(a.RoomId, conn.RemoteAddr().String())

	c.RoomId = ""
	c.UserId = a.UserId

	//userDevice, _ := service.GetUserDevice(ctx, c.ClientName())
	//service.DelUserDeviceInRoom(ctx, c.ClientName())
	//
	//service.DelRoomUser(ctx, userDevice.RoomId, c.ClientName())

	//noticeLeaveRoom(ctx, conn, packet.Header.MessageId, a.UserId, a.RoomId)

	return protocol.NewResponseOK(packet, nil), nil
}

func changeRoom(ctx context.Context, c net.Conn, packet *protocol.Packet, body *protocol.MessageBodySignalChangeRoom, cc *thread.ConnectClient) (*protocol.Packet, error) {

	if util.IsEmpty(body.NewRoomId) {
		return protocol.NewResponseError(packet, err.InvalidRequest.Format("NewRoomId")), nil
	}

	thread.RemRoomChannel(body.OldRoomId, c.RemoteAddr().String())
	thread.SetRoomChannels(body.NewRoomId, c.RemoteAddr().String())

	cc.RoomId = body.NewRoomId
	cc.UserId = body.UserId

	//info, _ := service.GetUserDevice(ctx, c.ClientName())
	//service.DelRoomUser(ctx, info.RoomId, c.ClientName())
	//service.SetRoomUser(ctx, body.RoomId, c.ClientName())
	//
	//service.SetUserDevice2InRoom(ctx, c.ClientName(), body.RoomId)

	//noticeLeaveRoom(ctx, c, packet.Header.MessageId, body.UserId, body.OldRoomId)
	//noticeJoinRoom(ctx, c, packet.Header.MessageId, body.UserId, body.NewRoomId)

	//body.RoomBlocked = service.GetRoomBlocked(ctx, body.NewRoomId)
	//body.Blocked = service.GetRoomMemberBlocked(ctx, body.OldRoomId, body.NewRoomId)

	return protocol.NewResponseOK(packet, body), nil
}
