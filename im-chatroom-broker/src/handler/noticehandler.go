package handler

import (
	"im-chatroom-broker/protocol"
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

func noticeJoinRoom(msgId, userId, roomId string) {

	//if util.IsEmpty(roomId) {
	//	return
	//}
	//
	//packet := buildNoticePacket(msgId, roomId, userId, protocol.TypeNoticeJoinRoom)
	//
	//deliver.Deliver2Worker(packet)
	//Deliver2AnotherBroker(packet)
}

func noticeLeaveRoom(msgId, userId, roomId string) {

	//if util.IsEmpty(roomId) {
	//	return
	//}
	//
	//packet := buildNoticePacket(msgId, roomId, userId, protocol.TypeNoticeLeaveRoom)
	//
	//deliver.Deliver2Worker(packet)
	//Deliver2AnotherBroker(packet)
}

func buildNoticePacket(msgId, roomId string, userId string, noticeType uint32) *protocol.Packet {
	packet := &protocol.Packet{
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
