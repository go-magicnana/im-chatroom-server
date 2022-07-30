package handler

import (
	"golang.org/x/net/context"
	err "im-chatroom-broker/error"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/thread"
	"im-chatroom-broker/zaplog"
	"net"
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

func (d DefaultHandler) Handle(ctx context.Context, conn net.Conn, packet *protocol.Packet,c *thread.ConnectClient) (*protocol.Packet, error) {
	ret := protocol.NewResponseError(packet, err.CommandNotAllow)
	switch packet.Header.Type {
	case protocol.TypeDefaultHeartBeat:
		a := protocol.JsonDefaultHearBeat(packet.Body)
		packet.Body = a
		return heartbeat(ctx, conn, packet)
	}
	return ret, nil
}

func heartbeat(ctx context.Context, conn net.Conn, packet *protocol.Packet) (*protocol.Packet, error) {

	body := packet.Body.(*protocol.MessageBodyDefaultHeartBeat)

	if body.Password == protocol.TypeDefaultHeartBeatPassword {

		cs := 0

		thread.RanChannel(func(key, value any) bool {
			cs++

			cc := value.(*thread.ConnectClient)
			zaplog.Logger.Debugf("ThreadContext HeartBeat %s %v %d %d", key, cc,cs,thread.Count.Load())

			return true
		})

		return protocol.NewResponseOK(packet, "OK"), nil
	} else {
		return protocol.NewResponseOK(packet, "QUIT"), nil

	}
}
