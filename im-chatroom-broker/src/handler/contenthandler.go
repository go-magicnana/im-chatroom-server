package handler

import (
	"golang.org/x/net/context"
	"golang.org/x/time/rate"
	"im-chatroom-broker/ctx"
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

var _limit *rate.Limiter

func SingleContentHandler() *ContentHandler {
	onceContentHandler.Do(func() {
		contentHandler = &ContentHandler{}
		_limit = rate.NewLimiter(100,100)
	})

	return contentHandler
}

type ContentHandler struct{}

/**
TypeContentText  = 4101
TypeContentEmoji = 4102
TypeContentReply = 4103
*/

func (d ContentHandler) Handle(c *ctx.Context, packet *protocol.Packet) (*protocol.Packet, error) {

	zaplog.Logger.Debugf("Handler content %s", c.ClientName)

	ret := protocol.NewResponseError(packet, err.TypeNotAllow)

	switch packet.Header.Type {
	case protocol.TypeContentText:
		a := protocol.JsonContentText(packet.Body)
		packet.Body = a
		return text(c, packet)

	case protocol.TypeContentEmoji:
		a := protocol.JsonContentText(packet.Body)
		packet.Body = a
		return text(c, packet)

	case protocol.TypeContentAt:
		a := protocol.JsonContentAt(packet.Body)
		packet.Body = a

		return at(c, packet, a)

	case protocol.TypeContentReply:
		a := protocol.JsonContentReply(packet.Body)
		packet.Body = a

		return reply(c, packet, a)

	}
	zaplog.Logger.Debugf("Handler content %s", c.ClientName)

	return ret, nil
}

func todeliver(c *ctx.Context, packet *protocol.Packet) (*protocol.Packet, error) {
	userId := packet.Header.From.UserId
	user := service.GetUserInfo(userId)
	if user != nil {
		packet.Header.From = *user
	}

	if packet.Header.Target == protocol.TargetRoom {
		_limit.Wait(context.Background())
	}

	packet.Header.Flow = protocol.FlowDeliver
	packet.Header.Code = err.OK.Code
	packet.Header.Message = err.OK.Message

	mq.Deliver2Worker(true, packet)

	return protocol.NewResponseOK(packet, nil), nil
}

func text(c *ctx.Context, packet *protocol.Packet) (*protocol.Packet, error) {

	if util.IsEmpty(packet.Body.(*protocol.MessageBodyContentText).Content) {
		return nil, nil
	}

	return todeliver(c, packet)
}

func at(c *ctx.Context, packet *protocol.Packet, body *protocol.MessageBodyContentAt) (*protocol.Packet, error) {

	user := service.GetUserInfo(body.AtUser.UserId)
	if user != nil {
		body.AtUser = *user
	}
	packet.Body = body

	return todeliver(c, packet)
}

func reply(c *ctx.Context, packet *protocol.Packet, body *protocol.MessageBodyContentReply) (*protocol.Packet, error) {

	user := service.GetUserInfo(body.ReplyUser.UserId)
	if user != nil {
		body.ReplyUser = *user
	}

	packet.Body = body

	return todeliver(c, packet)

}
