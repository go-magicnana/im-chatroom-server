package handler

import (
	"golang.org/x/net/context"
	context2 "im-chatroom-broker/context"
	err "im-chatroom-broker/error"
	"im-chatroom-broker/mq"
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

		return at(ctx, c, packet)

	case protocol.TypeContentReply:
		a := protocol.JsonContentReply(packet.Body)
		packet.Body = a

		return reply(ctx, c, packet)

	}
	return ret, nil
}

func deliver(ctx context.Context,c *context2.Context,packet *protocol.Packet)(*protocol.Packet, error){
	user, e1 := GetUserInfo(ctx, c.UserKey())
	if e1 != nil {
		return nil, e1
	}

	packet.Header.From = *user
	packet.Header.Flow = protocol.FlowDeliver

	if packet.Header.Target == protocol.TargetRoom {
		mq.DeliverMessageToRoom(ctx, c, packet)
	} else {
		mq.DeliverMessageToUser(ctx, c, packet)
	}

	return protocol.NewResponseOK(packet,nil), nil
}

func text(ctx context.Context, c *context2.Context, packet *protocol.Packet) (*protocol.Packet, error) {

	if util.IsEmpty(packet.Body.(protocol.MessageBodyContentText).Content){
		return nil,nil
	}

	return deliver(ctx,c,packet)
}

func at(ctx context.Context, c *context2.Context, packet *protocol.Packet, body *protocol.MessageBodyContentAt) (*protocol.Packet, error) {





	user, e1 := GetUserInfo(ctx, c.UserKey())
	if e1 != nil {
		return nil, e1
	}

	atUser, e2 := GetUserInfo(ctx, body.AtUserKey)
	if e2 != nil {
		return nil, e2
	}

	body.AtUserId = atUser.UserId
	body.AtUserName = atUser.Name
	body.AtUserAvatar = atUser.Avatar

	packet.Header.From = *user
	packet.Header.Flow = protocol.FlowDeliver

	if packet.Header.Target == protocol.TargetRoom {
		mq.DeliverMessageToRoom(ctx, c, packet)
	} else {
		mq.DeliverMessageToUser(ctx, c, packet)
	}

	return nil, nil
}

func reply(ctx context.Context, c *context2.Context, packet *protocol.Packet, body *protocol.MessageBodyContentReply) (*protocol.Packet, error) {

	user, e1 := GetUserInfo(ctx, c.UserKey())
	if e1 != nil {
		return nil, e1
	}

	atUser, e2 := GetUserInfo(ctx, body.ReplyUserKey)
	if e2 != nil {
		return nil, e2
	}

	body.ReplyUserId = atUser.UserId
	body.ReplyUserName = atUser.Name
	body.ReplyUserAvatar = atUser.Avatar

	packet.Header.From = *user
	packet.Header.Flow = protocol.FlowDeliver

	if packet.Header.Target == protocol.TargetRoom {
		mq.DeliverMessageToRoom(ctx, c, packet)
	} else {
		mq.DeliverMessageToUser(ctx, c, packet)
	}

	return nil, nil
}
