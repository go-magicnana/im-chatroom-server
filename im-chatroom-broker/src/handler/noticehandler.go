package handler

import (
	"golang.org/x/net/context"
	context2 "im-chatroom-broker/context"
	err "im-chatroom-broker/error"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/util"
	"sync"
)

var onceNoticeHandler sync.Once

var noticeHandler *NoticeHandler

func SingleNoticeHandler() *NoticeHandler {
	onceNoticeHandler.Do(func() {
		noticeHandler = &NoticeHandler{}
	})

	return noticeHandler
}

type NoticeHandler struct{}

func (d NoticeHandler) Handle(ctx context.Context, c *context2.Context, packet *protocol.Packet) (*protocol.Packet, error) {
	ret := protocol.NewResponseError(packet, err.CommandNotAllow)

	switch packet.Header.Type {
	case protocol.TypeSignalJoinRoom:
		a := protocol.JsonSignalJoinRoom(packet.Body)
		return noticeJoinRoom(ctx, c, packet, a)

		//case protocol.TypeSignalLeaveRoom:
		//	return leaveRoom(ctx, c, packet)
		//
		//case protocol.TypeSignalChangeRoom:
		//	a := protocol.JsonSignalChangeRoom(packet.Body)
		//	return changeRoom(ctx, c, packet, a)
	}

	return ret, nil
}

func noticeJoinRoom(ctx context.Context, c *context2.Context, packet *protocol.Packet, body *protocol.MessageBodySignalJoinRoom) (*protocol.Packet, error) {

	if util.IsEmpty(body.RoomId) {
		return protocol.NewResponseError(packet, err.InvalidRequest.Format("roomId")), nil
	}

	//server.Write(ctx, c, packet)

	return protocol.NewResponseOK(packet, nil), nil
}
