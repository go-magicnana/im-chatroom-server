package mq

import (
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"golang.org/x/net/context"
	"im-chatroom-broker/config"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/thread"
	"im-chatroom-broker/util"
	"im-chatroom-broker/zaplog"
	"strings"
)

const (
	RoomGroup = "imchatroom_room_group"
	RoomTopic = "imchatroom_room_topic"

	OneGroup = "imchatroom_one_group_"
	OneTopic = "imchatroom_one_topic_"
)

var MyName = ""

var _producer rocketmq.Producer
var _consumer1 rocketmq.PushConsumer
var _consumer2 rocketmq.PushConsumer

var PushRoomChannel chan *protocol.PacketMessage
var Push2OneChannel chan *protocol.PacketMessage

func init() {

	if util.IsNotEmpty(config.OP.Ip) {
		MyName = broker2name(config.OP.Ip + ":" + config.OP.Port)

	} else {
		ip := util.GetBrokerIp()
		MyName = broker2name(ip + ":" + config.OP.Port)

	}

	PushRoomChannel = make(chan *protocol.PacketMessage, 500)
	Push2OneChannel = make(chan *protocol.PacketMessage, 1500)
	_producer = newProducer()
	_consumer1 = newConsumerRoom()
	_consumer2 = newConsumerOne()

	c, _ := context.WithCancel(context.Background())

	go SendSync2Room(c)
	go SendSync2One(c)

}

func newProducer() rocketmq.Producer {
	zaplog.Logger.Infof("Init RocketMQ Producer %s", config.OP.RocketMQ.Address)
	p, _ := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{config.OP.RocketMQ.Address})),
		producer.WithQueueSelector(producer.NewHashQueueSelector()),
		producer.WithRetry(2),
	)
	err := p.Start()
	if err != nil {
		util.Panic(err)
	}

	return p
}

func newConsumerRoom() rocketmq.PushConsumer {

	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName(RoomGroup),
		consumer.WithConsumerModel(consumer.BroadCasting),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{config.OP.RocketMQ.Address})),
	)
	err := c.Subscribe(RoomTopic, consumer.MessageSelector{},
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

			for i := range msgs {

				p := &protocol.PacketMessage{}
				json.Unmarshal(msgs[i].Body, p)

				zaplog.Logger.Debugf("CsumRoom %s %s C:%d T:%d F:%d %v %v", RoomTopic, p.Packet.Header.MessageId, p.Packet.Header.Command, p.Packet.Header.Type, p.Packet.Header.Flow, p.Packet.Body, msgs[i].MsgId)

				if p.Packet.Header.Target == protocol.TargetRoom {

					roomClients := thread.GetRoomChannels(p.Packet.Header.To)

					for _, clientName := range roomClients {
						cc := thread.GetChannel(clientName.(string))

						if cc != nil && cc.Channel != nil {
							cc.Channel <- protocol.NewResponse(p.Packet)
						}
					}

				}

			}
			zaplog.Logger.Debugf("---------------------------------------------------------")

			return consumer.ConsumeSuccess, nil
		})

	if err != nil {
		util.Panic(err)
	}
	// Note: start after subscribe
	err = c.Start()
	if err != nil {
		util.Panic(err)
	}

	zaplog.Logger.Infof("Init RocketMQ Consumer-Room %s", config.OP.RocketMQ.Address)

	return c
}

//
//func write2queue(queue *context2.Queue,packet *protocol.Packet){
//	/*
//
//	var next *list.Element
//	for e := l.Front(); e != nil; e = next {
//	    next = e.Next()
//	    l.Remove(e)
//	}
//
//	for {
//	    e := l.Front()
//	    if e == nil {
//	            break
//	    }
//	    l.Remove(e)
//	}
//
//	 */
//
//
//	for {
//		e := queue.Front()
//		if e == nil {
//			break
//		}
//
//		c,b:=service.GetUserContext(e.(string))
//		if !b || c==nil {
//			continue
//		}
//
//		serializer.SingleJsonSerializer().Write(c, packet)
//	}
//}

func newConsumerOne() rocketmq.PushConsumer {

	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName(OneGroup+MyName),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{config.OP.RocketMQ.Address})),
	)
	err := c.Subscribe(OneTopic+MyName, consumer.MessageSelector{},
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

			for i := range msgs {
				p := &protocol.PacketMessage{}
				json.Unmarshal(msgs[i].Body, p)

				zaplog.Logger.Debugf("ConsumeOne %s %s %v", OneTopic+MyName, msgs[i].MsgId, p)

				c := thread.GetChannel(p.ClientName)

				if c != nil && c.Channel != nil {
					c.Channel <- protocol.NewResponse(p.Packet)

				}

			}
			return consumer.ConsumeSuccess, nil

		})
	if err != nil {
		util.Panic(err)
	}
	// Note: start after subscribe
	err = c.Start()
	if err != nil {
		util.Panic(err)
	}

	zaplog.Logger.Infof("Init RocketMQ Consumer-One %s", config.OP.RocketMQ.Address)

	return c
}

func createTopic(topicName string) {
	endPoint := []string{config.OP.RocketMQ.Address}
	// 创建主题
	testAdmin, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver(endPoint)))
	if err != nil {
		fmt.Printf("connection error: %s\n", err.Error())
	}
	err = testAdmin.CreateTopic(context.Background(), admin.WithTopicCreate(topicName))
	if err != nil {
		fmt.Printf("createTopic error: %s\n", err.Error())
	}
}

func HashString(s string) int {
	val := []byte(s)
	var h int32

	for idx := range val {
		h = 31*h + int32(val[idx])
	}

	return int(h)
}

func sendSync(topic, key string, message []byte) (*primitive.SendResult, error) {
	msg := &primitive.Message{
		Topic: topic,
		Body:  message,
	}
	msg.WithProperty(primitive.PropertyShardingKey, key)
	res, err := _producer.SendSync(context.Background(), msg)

	return res, err
}

func SendSync2One(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			return
		case p := <-Push2OneChannel:
			topic := OneTopic + broker2name(p.Broker)
			msg, _ := json.Marshal(p)
			_, err := sendSync(topic, p.ClientName, msg)
			zaplog.Logger.Debugf("SendSync %s %s C:%d T:%d F:%d %v %v", topic, p.Packet.Header.MessageId, p.Packet.Header.Command, p.Packet.Header.Type, p.Packet.Header.Flow, p.Packet.Body, err)
		default:
			continue
		}
	}

}

func SendSync2Room(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case p := <-PushRoomChannel:
			msg, _ := json.Marshal(p)
			ret, _ := sendSync(RoomTopic, p.ClientName, msg)
			zaplog.Logger.Debugf("SendSync %s %s C:%d T:%d F:%d %v %v", RoomTopic, p.Packet.Header.MessageId, p.Packet.Header.Command, p.Packet.Header.Type, p.Packet.Header.Flow, p.Packet.Body, ret.MsgID)
		default:
			continue
		}
	}
}

func broker2name(broker string) string {
	broker = strings.ReplaceAll(broker, ".", "_")
	broker = strings.ReplaceAll(broker, ":", "_")
	return broker
}
