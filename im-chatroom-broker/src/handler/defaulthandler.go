package handler

import (
	"golang.org/x/net/context"
	context2 "im-chatroom-broker/context"
	err "im-chatroom-broker/error"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/service"
	"strconv"
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

	body := packet.Body.(*protocol.MessageBodyDefaultHeartBeat)

	if body.Password == protocol.TypeDefaultHeartBeatPassword {


		cs := 0
		rs := 0

		service.RangeUserContextAll(func(key, value any) bool {
			cs ++
			return true
		})


		service.RangeRoom("1", func(key, value any) bool {
			rs ++
			return true
		})

		return protocol.NewResponseOK(packet, "OK "+strconv.Itoa(cs)+" room1"+strconv.Itoa(rs)), nil
	} else {
		service.Close(ctx, c,nil)
		return nil, nil
	}
}
