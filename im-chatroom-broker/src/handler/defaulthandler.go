package handler

import (
	"fmt"
	"im-chatroom-broker/ctx"
	err "im-chatroom-broker/error"
	"im-chatroom-broker/protocol"
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

func (d DefaultHandler) Handle(c *ctx.Context, packet *protocol.Packet) (*protocol.Packet, error) {
	ret := protocol.NewResponseError(packet, err.CommandNotAllow)
	switch packet.Header.Type {
	case protocol.TypeDefaultHeartBeat:
		a := protocol.JsonDefaultHearBeat(packet.Body)
		packet.Body = a
		return heartbeat(c, packet)
	}
	return ret, nil
}

func heartbeat(c *ctx.Context, packet *protocol.Packet) (*protocol.Packet, error) {

	body := packet.Body.(*protocol.MessageBodyDefaultHeartBeat)

	if body.Password == protocol.TypeDefaultHeartBeatPassword {

		size := ctx.ConnCount()

		return protocol.NewResponseOK(packet, "OK "+fmt.Sprintf("%d", size)), nil
	} else {
		return protocol.NewResponseOK(packet, "QUIT"), nil

	}
}
