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
	"im-chatroom-broker/serializer"
	"im-chatroom-broker/service"
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

func init() {

	ip := util.GetBrokerIp()
	MyName = broker2name(ip + ":" + config.OP.Port)

	_producer = newProducer()
	_consumer1 = newConsumerRoom()
	//_consumer2 = newConsumerOne()
}

func newProducer() rocketmq.Producer {
	zaplog.Logger.Infof("Init RocketMQ Producer %s", config.OP.RocketMQ.Address)
	p, _ := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{config.OP.RocketMQ.Address})),
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
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{config.OP.RocketMQ.Address})),
	)
	err := c.Subscribe(RoomTopic, consumer.MessageSelector{},
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

			for i := range msgs {

				p := &protocol.Packet{}
				json.Unmarshal(msgs[i].Body, p)

				zaplog.Logger.Debugf("CsumRoom %s %s C:%d T:%d F:%d %v %v", RoomTopic, p.Header.MessageId, p.Header.Command, p.Header.Type, p.Header.Flow, p.Body, msgs[i].MsgId)

				if p.Header.Target == protocol.TargetRoom {
					go service.RangeRoom(p.Header.To, func(key, value any) bool {

						c, b := service.GetUserContext(key.(string))
						if b && c != nil {
							serializer.SingleJsonSerializer().Write(c, p)
						}
						return true

					})

					//b, e := service.GetRoomMembers(ctx, p.Header.To)
					//
					//if e == nil {
					//	for _, v := range b {
					//
					//		broker, _ := service.GetUserDeviceBroker(ctx, v)
					//
					//		if util.IsEmpty(broker) {
					//			continue
					//		}
					//
					//		m := protocol.PacketMessage{
					//			ClientName: v,
					//			Packet:     *p,
					//		}
					//
					//		SendSync2One(broker, &m)
					//
					//	}
					//}

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

				c, e := service.GetUserContext(p.ClientName)

				if e {
					serializer.SingleJsonSerializer().Write(c, &p.Packet)

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

func sendSync(topic string, message []byte) (*primitive.SendResult, error) {
	msg := &primitive.Message{
		Topic: topic,
		Body:  message,
	}
	res, err := _producer.SendSync(context.Background(), msg)

	return res, err
}

func SendSync2One(broker string, packet *protocol.PacketMessage) {

	topic := OneTopic + broker2name(broker)
	msg, _ := json.Marshal(packet)
	_, err := sendSync(topic, msg)
	p := packet.Packet
	zaplog.Logger.Debugf("SendSync %s %s C:%d T:%d F:%d %v %v", topic, p.Header.MessageId, p.Header.Command, p.Header.Type, p.Header.Flow, p.Body, err)
}

func SendSync2Room(p *protocol.Packet) string {

	msg, _ := json.Marshal(p)
	ret, _ := sendSync(RoomTopic, msg)
	zaplog.Logger.Debugf("SendSync %s %s C:%d T:%d F:%d %v %v", RoomTopic, p.Header.MessageId, p.Header.Command, p.Header.Type, p.Header.Flow, p.Body, ret.MsgID)
	return ret.MsgID
}

func broker2name(broker string) string {
	broker = strings.ReplaceAll(broker, ".", "_")
	broker = strings.ReplaceAll(broker, ":", "_")
	return broker
}
