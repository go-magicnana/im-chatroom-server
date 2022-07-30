package handler

import (
	"golang.org/x/net/context"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/thread"
	"net"
)

type Handler interface {
	Handle(ctx context.Context,conn net.Conn, packet *protocol.Packet,c *thread.ConnectClient) (*protocol.Packet, error)
}
