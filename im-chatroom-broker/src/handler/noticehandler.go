package handler

import (
	"golang.org/x/net/context"
	context2 "im-chatroom-broker/context"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/util"
)

//var onceNoticeHandler sync.Once
//
//var noticeHandler *NoticeHandler

//func SingleNoticeHandler() *NoticeHandler {
//	onceNoticeHandler.Do(func() {
//		noticeHandler = &NoticeHandler{}
//	})
//
//	return noticeHandler
//}

//type NoticeHandler struct{}

//func (d NoticeHandler) Handle(ctx context.Context, c *context2.Context, packet *protocol.Packet) {
//
//	switch packet.Header.Type {
//	case protocol.TypeSignalJoinRoom:
//		a := protocol.MessageBodyNoticeJoinRoom{}
//		noticeJoinRoom(ctx, c, packet, a)
//
//		//case protocol.TypeSignalLeaveRoom:
//		//	return leaveRoom(ctx, c, packet)
//		//
//		//case protocol.TypeSignalChangeRoom:
//		//	a := protocol.JsonSignalChangeRoom(packet.Body)
//		//	return changeRoom(ctx, c, packet, a)
//	}
//
//}

func noticeJoinRoom(ctx context.Context, c *context2.Context, msgId,roomId string) {

	if util.IsEmpty(roomId) {
		return
	}

	packet := buildNoticePacket(msgId,roomId, c.UserId(), protocol.TypeNoticeJoinRoom)

	deliver(ctx, c, &packet)
}

func noticeLeaveRoom(ctx context.Context, c *context2.Context, msgId,roomId string) {

	if util.IsEmpty(roomId) {
		return
	}

	packet := buildNoticePacket(msgId,roomId, c.UserId(), protocol.TypeNoticeLeaveRoom)

	deliver(ctx, c, &packet)
}

func buildNoticePacket(msgId,roomId string, userId string, noticeType uint32) protocol.Packet {
	packet := protocol.Packet{
		Header: protocol.MessageHeader{
			MessageId: msgId,
			Command:   protocol.CommandNotice,
			Target:    protocol.TargetRoom,
			To:        roomId,
			Type:      noticeType,
		},
		Body: protocol.MessageBodyNoticeJoinRoom{
			RoomId: roomId,
			UserId: userId,
		},
	}
	return packet
}
