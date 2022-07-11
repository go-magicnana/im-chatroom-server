package serializer

import (
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/server"
)

type Serializer interface {
	Enpack(packet *protocol.Packet,c *server.Context) ([]byte,error)
	Depack(packet *protocol.Packet,c *server.Context)	(*protocol.Packet,error)
}
