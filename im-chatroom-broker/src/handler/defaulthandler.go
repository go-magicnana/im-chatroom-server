package handler

import (
	"go.uber.org/atomic"
	"im-chatroom-broker/ctx"
	err "im-chatroom-broker/error"
	"im-chatroom-broker/protocol"
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

var Connections *atomic.Int64 = atomic.NewInt64(0)

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

		return protocol.NewResponseOK(packet, "OK "+strconv.Itoa(int(Connections.Load()))), nil
	} else {
		return protocol.NewResponseOK(packet, "QUIT"), nil

	}
}
