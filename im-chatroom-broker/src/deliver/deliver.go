package deliver

import (
	"im-chatroom-broker/ctx"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/serializer"
	"im-chatroom-broker/service"
)

var _serializer *serializer.JsonSerializer

func init() {
	_serializer = serializer.NewJsonSerializer()
}

func Deliver2Broker(broker string, packet *protocol.Packet) {

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
