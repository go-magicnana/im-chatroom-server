package handler

import (
	"golang.org/x/net/context"
	err "im-chatroom-broker/error"
	"im-chatroom-broker/mq"
	"im-chatroom-broker/service"
	"im-chatroom-broker/thread"
	"im-chatroom-broker/zaplog"
	"math"
	"net"

	//"im-chatroom-broker/mq"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/util"
	"sync"
)

var onceContentHandler sync.Once

var contentHandler *ContentHandler

func SingleContentHandler() *ContentHandler {
	onceContentHandler.Do(func() {
		contentHandler = &ContentHandler{}
	})

	return contentHandler
}

type ContentHandler struct{}

/**
TypeContentText  = 4101
TypeContentEmoji = 4102
TypeContentReply = 4103
*/

func (d ContentHandler) Handle(ctx context.Context, conn net.Conn, packet *protocol.Packet, c *thread.ConnectClient) (*protocol.Packet, error) {
	ret := protocol.NewResponseError(packet, err.TypeNotAllow)

	switch packet.Header.Type {
	case protocol.TypeContentText:
		a := protocol.JsonContentText(packet.Body)
		packet.Body = a
		return text(ctx, conn, packet, c)

	case protocol.TypeContentEmoji:
		a := protocol.JsonContentText(packet.Body)
		packet.Body = a
		return text(ctx, conn, packet, c)

	case protocol.TypeContentAt:
		a := protocol.JsonContentAt(packet.Body)
		packet.Body = a

		return at(ctx, conn, packet, a, c)

	case protocol.TypeContentReply:
		a := protocol.JsonContentReply(packet.Body)
		packet.Body = a

		return reply(ctx, conn, packet, a, c)

	}
	return ret, nil
}

func deliver(ctx context.Context, conn net.Conn, packet *protocol.Packet, c *thread.ConnectClient) (*protocol.Packet, error) {
	userId := packet.Header.From.UserId
	user := service.GetUserInfo(ctx, userId)
	if user != nil {
		packet.Header.From = *user
	}

	packet.Header.Flow = protocol.FlowDeliver
	packet.Header.Code = err.OK.Code
	packet.Header.Message = err.OK.Message

	if packet.Header.Target == protocol.TargetRoom {
		zaplog.Logger.Debugf("Deliver %s C:%d T:%d F:%d %v", packet.Header.MessageId, packet.Header.Command, packet.Header.Type, packet.Header.Flow, packet.Body)

		//deliver2ConsumerRoom(c, conn, packet)

		go deliver2BrokerRoom(packet, conn, c)

	} else {

		//ret := service.GetUserClients(ctx, packet.Header.To)
		//
		//for _, v := range ret {
		//
		//	if v == c.ClientName() {
		//		continue
		//	}
		//
		//	msg := &protocol.PacketMessage{
		//		ClientName: v,
		//		Packet:  *packet,
		//	}
		//
		//	broker, _ := service.GetUserDeviceBroker(ctx, v)
		//
		//	mq.SendSync2One(broker, msg)
		//	//fmt.Println(msg,broker)
		//}

	}

	return protocol.NewResponseOK(packet, nil), nil
}

func deliver2ConsumerRoom(c *thread.ConnectClient, conn net.Conn, packet *protocol.Packet) {
	pm := &protocol.PacketMessage{
		Broker:     c.Broker,
		ClientName: conn.RemoteAddr().String(),
		Packet:     packet,
	}
	mq.PushRoomChannel <- pm
}

func deliver2BrokerRoom(packet *protocol.Packet, conn net.Conn, c *thread.ConnectClient) {

	defer func() {
		zaplog.Logger.Errorf("DeliverRecover1 %s DropPacket %v", conn.RemoteAddr(), recover())
	}()

	cs := thread.GetRoomChannels(packet.Header.To)
	if cs != nil {

		deliver2BrokerRoomBatch(packet, cs,conn)
		//for _, v := range cs {
		//	clientName := v.(string)
		//
		//	cc := thread.GetChannel(clientName)
		//
		//	if cc != nil && cc.Channel != nil {
		//		cc.Channel <- protocol.NewResponse(packet)
		//
		//	}
		//}
	}
}

func deliver2BrokerRoomBatch(packet *protocol.Packet, cs []interface{},conn net.Conn) {
	iss := spiltList(cs, 100)

	for _, v := range iss {

		if v != nil {
			go deliver2BrokerRoomBatchInner(v, packet,conn)
		}

	}

}

func deliver2BrokerRoomBatchInner(clients []interface{}, packet *protocol.Packet,conn net.Conn) {

	defer func() {
		zaplog.Logger.Errorf("DeliverRecover1 %s DropPacket %v", conn.RemoteAddr(), recover())
	}()

	if clients != nil {
		for _, client := range clients {
			clientName := client.(string)
			cc := thread.GetChannel(clientName)

			if cc != nil && cc.Channel != nil {
				cc.Channel <- protocol.NewResponse(packet)

			}
		}
	}
}

func spiltList(list []interface{}, size int) [][]interface{} {
	lens := len(list)
	mod := math.Ceil(float64(lens) / float64(size))
	spiltList := make([][]interface{}, 0)
	for i := 0; i < int(mod); i++ {
		tmpList := make([]interface{}, 0, size)
		if i == int(mod)-1 {
			tmpList = list[i*size:]
		} else {
			tmpList = list[i*size : i*size+size]
		}
		spiltList = append(spiltList, tmpList)
	}
	return spiltList
}

func text(ctx context.Context, conn net.Conn, packet *protocol.Packet, c *thread.ConnectClient) (*protocol.Packet, error) {

	if util.IsEmpty(packet.Body.(*protocol.MessageBodyContentText).Content) {
		return nil, nil
	}

	return deliver(ctx, conn, packet, c)
}

func at(ctx context.Context, conn net.Conn, packet *protocol.Packet, body *protocol.MessageBodyContentAt, c *thread.ConnectClient) (*protocol.Packet, error) {

	user := service.GetUserInfo(ctx, body.AtUser.UserId)
	if user != nil {
		body.AtUser = *user
	}
	packet.Body = body

	return deliver(ctx, conn, packet, c)
}

func reply(ctx context.Context, conn net.Conn, packet *protocol.Packet, body *protocol.MessageBodyContentReply, c *thread.ConnectClient) (*protocol.Packet, error) {

	user := service.GetUserInfo(ctx, body.ReplyUser.UserId)
	if user != nil {
		body.ReplyUser = *user
	}

	packet.Body = body

	return deliver(ctx, conn, packet, c)

}
