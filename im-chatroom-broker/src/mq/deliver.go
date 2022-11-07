package mq

import (
	"im-chatroom-broker/ctx"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/serializer"
	"im-chatroom-broker/service"
	"im-chatroom-broker/zaplog"
	"runtime"
	"strings"
)

var _serializer *serializer.JsonSerializer
var _worker *pool

func init() {
	_serializer = serializer.NewJsonSerializer()
	_worker = newPool(runtime.NumCPU(), 65536)
	_worker.start()
}

func deliver(toMQ bool, packet *protocol.Packet) {

	zaplog.Logger.Debugf("Worker deliver %d %d %d %d %s", packet.Header.Command, packet.Header.Type, packet.Header.Target, packet.Header.Flow, packet.Header.To)

	if packet.Header.Target == protocol.TargetRoom {
		deliver2ThisBrokerRoom(packet)

		if toMQ {
			deliver2AnotherBrokerRoom(packet)
		}

	} else {
		deliver2User(toMQ, packet)
	}
}

func deliver2ThisBrokerRoom(packet *protocol.Packet) {
	if packet.Header.Target == protocol.TargetRoom {
		cs := service.GetRoomClients(packet.Header.To)
		if cs != nil {

			for _, v := range cs {

				cc := ctx.GetContext(v.(string))

				if cc != nil && cc.Conn != nil {

					retBuf, _ := _serializer.EncodePacket(packet)
					buf, _ := _serializer.Encode(retBuf)
					cc.Conn.AsyncWrite(buf, nil)
				}
			}
		}
	}
}

func deliver2AnotherBrokerRoom(packet *protocol.Packet) {
	if packet.Header.Target == protocol.TargetRoom {
		cs := service.GetBrokerInstances()
		if cs != nil {

			for _, v := range cs {

				if v == ctx.BrokerAddress {
					continue
				}

				exist := service.GetBrokerRoomExist(v, packet.Header.To)

				if exist {
					Deliver2MQ(v, packet)
				}
			}
		}
	}
}

func deliver2User(toMQ bool, packet *protocol.Packet) {
	clients := service.GetUserClients(packet.Header.To)
	if clients != nil {
		for _, v := range clients {
			bc := strings.Split(v, "-")
			if bc[0] == ctx.BrokerAddress {

				cc := ctx.GetContext(bc[1])

				if cc != nil {

					retBuf, _ := _serializer.EncodePacket(packet)
					buf, _ := _serializer.Encode(retBuf)
					cc.Conn.AsyncWrite(buf, nil)
				}

			} else {
				if toMQ {
					Deliver2MQ(bc[0], packet)
				}
			}
		}
	}
}

func Deliver2Worker(toMQ bool, packet *protocol.Packet) {
	_worker.addTask(&protocol.PacketDeliver{Packet: packet, ToMQ: toMQ})
	zaplog.Logger.Debugf("Worker addTask %d %d %d %d %s", packet.Header.Command, packet.Header.Type, packet.Header.Target, packet.Header.Flow, packet.Header.To)

}
