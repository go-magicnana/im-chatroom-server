package handler

import (
	"golang.org/x/net/context"
	context2 "im-chatroom-broker/context"
	err "im-chatroom-broker/error"
	"im-chatroom-broker/mq"
	"im-chatroom-broker/service"
	"im-chatroom-broker/zaplog"

	//"im-chatroom-broker/mq"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/util"
	"sync"
)

var onceContentHandler sync.Once

var contentHandler *ContentHandler

func SingleContentHandler() *ContentHandler {
	onceDefaultHandler.Do(func() {
		contentHandler = &ContentHandler{}
	})

	return contentHandler
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
		a := protocol.JsonContentText(packet.Body)
		packet.Body = a
		return text(ctx, c, packet)

	case protocol.TypeContentEmoji:
		a := protocol.JsonContentText(packet.Body)
		packet.Body = a
		return text(ctx, c, packet)

	case protocol.TypeContentAt:
		a := protocol.JsonContentAt(packet.Body)
		packet.Body = a

		return at(ctx, c, packet, a)

	case protocol.TypeContentReply:
		a := protocol.JsonContentReply(packet.Body)
		packet.Body = a

		return reply(ctx, c, packet, a)

	}
	return ret, nil
}

func deliver(ctx context.Context, c *context2.Context, packet *protocol.Packet) (*protocol.Packet, error) {
	user, e1 := service.GetUserInfo(ctx, c.UserId())
	if e1 != nil {
		return nil, e1
	}

	packet.Header.From = *user
	packet.Header.Flow = protocol.FlowDeliver
	packet.Header.Code = err.OK.Code
	packet.Header.Message = err.OK.Message


	if packet.Header.Target == protocol.TargetRoom {
		zaplog.Logger.Debugf("Deliver RoomTopic %s C:%d T:%d F:%d %v", packet.Header.MessageId, packet.Header.Command, packet.Header.Type, packet.Header.Flow, packet.Body)

		mq.SendSync2Room(packet)
	} else {

		//ret := service.GetUserClients(ctx, packet.Header.To)
		//
		//for _, v := range ret {
		//
		//	if v == c.ClientName() {
		//		continue
		//	}
		//
		//	msg := &protocol.PacketMessage{
		//		ClientName: v,
		//		Packet:  *packet,
		//	}
		//
		//	broker, _ := service.GetUserDeviceBroker(ctx, v)
		//
		//	mq.SendSync2One(broker, msg)
		//	//fmt.Println(msg,broker)
		//}

	}

	return protocol.NewResponseOK(packet, nil), nil
}

func text(ctx context.Context, c *context2.Context, packet *protocol.Packet) (*protocol.Packet, error) {

	if util.IsEmpty(packet.Body.(*protocol.MessageBodyContentText).Content) {
		return nil, nil
	}

	return deliver(ctx, c, packet)
}

func at(ctx context.Context, c *context2.Context, packet *protocol.Packet, body *protocol.MessageBodyContentAt) (*protocol.Packet, error) {

	user, e := service.GetUserInfo(ctx, body.AtUser.UserId)
	if user == nil {
		return nil, e
	}

	body.AtUser = *user

	packet.Body = body

	return deliver(ctx, c, packet)
}

func reply(ctx context.Context, c *context2.Context, packet *protocol.Packet, body *protocol.MessageBodyContentReply) (*protocol.Packet, error) {

	user, e2 := service.GetUserInfo(ctx, body.ReplyUser.UserId)
	if user == nil {
		return nil, e2
	}

	body.ReplyUser = *user
	packet.Body = body

	return deliver(ctx, c, packet)

}
