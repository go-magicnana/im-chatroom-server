package deliver

import (
	"fmt"
	"im-chatroom-broker/ctx"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/serializer"
	"im-chatroom-broker/service"
	"im-chatroom-broker/zaplog"
	"runtime"
)

var _serializer *serializer.JsonSerializer
var _worker *pool

func init() {
	_serializer = serializer.NewJsonSerializer()
	_worker = newPool(runtime.NumCPU(), 65536)
	_worker.start()
}

func deliver(packet *protocol.Packet) {

	//mq.Deliver2MQ(broker,packet)

	zaplog.Logger.Debugf("Worker deliver %d %d %d %d %s", packet.Header.Command, packet.Header.Type, packet.Header.Target, packet.Header.Flow, packet.Header.To)

	if packet.Header.Target == protocol.TargetRoom {
		deliver2ThisBroker(packet)
	} else {
		fmt.Println("没写呢还")
	}
}

func deliver2ThisBroker(packet *protocol.Packet) {
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
	} else {
		fmt.Println("没写呢")
	}
}

func Deliver2Worker(packet *protocol.Packet) {

	_worker.addTask(packet)
	zaplog.Logger.Debugf("Worker addTask %d %d %d %d %s", packet.Header.Command, packet.Header.Type, packet.Header.Target, packet.Header.Flow, packet.Header.To)

}
