package deliver

import (
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
	_worker = newPool(runtime.NumCPU(),65536)
	_worker.start()
}

func Deliver2Broker(broker string, packet *protocol.Packet) {

	zaplog.Logger.Debugf("Worker deliver %s %d %d %d %d %s", broker,packet.Header.Command,packet.Header.Type,packet.Header.Target,packet.Header.Flow,packet.Header.To)

	if packet.Header.Target == protocol.TargetRoom {
		cs := service.GetRoomClients(broker, packet.Header.To)
		if cs != nil {

			for _, v := range cs {

				cc := ctx.GetContext(v)

				if cc != nil && cc.Conn != nil {
					retBuf, _ := _serializer.EncodePacket(packet)
					buf, _ := _serializer.Encode(retBuf)
					cc.Conn.AsyncWrite(buf, nil)
				}
			}
		}
	}

}

func Deliver2Worker(broker string,packet *protocol.Packet){
	_worker.addTask(&protocol.PacketMessage{
		Broker: broker,
		Packet: packet,
	})
	zaplog.Logger.Debugf("Worker addTask %s %d %d %d %d %s", broker,packet.Header.Command,packet.Header.Type,packet.Header.Target,packet.Header.Flow,packet.Header.To)

}
