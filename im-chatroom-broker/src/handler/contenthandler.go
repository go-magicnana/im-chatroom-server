package handler

import (
	"golang.org/x/net/context"
	context2 "im-chatroom-broker/context"
	err "im-chatroom-broker/error"
	"im-chatroom-broker/mq"
	"im-chatroom-broker/service"

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

	if packet.Header.Target == protocol.TargetRoom {
		mq.OneDeliver().ProduceRoom(packet)
	} else {

		ret := service.GetUserClients(ctx, packet.Header.To)

		for _, v := range ret {
			msg := &protocol.PacketMessage{
				UserKey: v,
				Packet:  *packet,
			}

			broker, _ := service.GetUserDeviceBroker(ctx, v)

			mq.OneDeliver().ProduceOne(broker, msg)
			//fmt.Println(msg,broker)
		}

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

	user, _ := service.GetUserInfo(ctx, body.AtUserId)

	body.AtUserId = user.UserId
	body.AtUserName = user.Name
	body.AtUserAvatar = user.Avatar

	return deliver(ctx, c, packet)
}

func reply(ctx context.Context, c *context2.Context, packet *protocol.Packet, body *protocol.MessageBodyContentReply) (*protocol.Packet, error) {

	atUser, e2 := service.GetUserInfo(ctx, body.ReplyUserId)
	if e2 != nil {
		return nil, e2
	}

	body.ReplyUserId = atUser.UserId
	body.ReplyUserName = atUser.Name
	body.ReplyUserAvatar = atUser.Avatar

	return deliver(ctx, c, packet)

}
