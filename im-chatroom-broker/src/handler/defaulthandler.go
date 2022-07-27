package handler

import (
	"golang.org/x/net/context"
	context2 "im-chatroom-broker/context"
	err "im-chatroom-broker/error"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/service"
	"sync"
)

var onceDefaultHandler sync.Once

var defaultHandler *DefaultHandler

func SingleDefaultHandler() *DefaultHandler {
	onceDefaultHandler.Do(func() {
		defaultHandler = &DefaultHandler{}
	})

	return defaultHandler
}

type DefaultHandler struct{}

func (d DefaultHandler) Handle(ctx context.Context, c *context2.Context, packet *protocol.Packet) (*protocol.Packet, error) {
	ret := protocol.NewResponseError(packet, err.CommandNotAllow)
	switch packet.Header.Type {
	case protocol.TypeDefaultHeartBeat:
		a := protocol.JsonDefaultHearBeat(packet.Body)
		packet.Body = a
		return heartbeat(ctx, c, packet)
	}
	return ret, nil
}

func heartbeat(ctx context.Context, c *context2.Context, packet *protocol.Packet) (*protocol.Packet, error) {

	body := packet.Body.(protocol.MessageBodyDefaultHeartBeat)

	if body.Password == protocol.TypeDefaultHeartBeatPassword {
		return protocol.NewResponseOK(packet, "OK"), nil
	} else {
		service.Close(ctx, c)
		return nil, nil
	}
}
