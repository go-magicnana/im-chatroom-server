package mq

import (
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"golang.org/x/net/context"
	"im-chatroom-gateway/protocol"
	"im-chatroom-gateway/util"
	"strings"
)

const (
	RoomGroup = "imchatroom_room_group"
	RoomTopic = "imchatroom_room_topic"

	OneGroup = "imchatroom_one_group_"
	OneTopic = "imchatroom_one_topic_"

	EndPoint = "192.168.3.242:9876"
)

var _producer rocketmq.Producer

//var _consumer1 rocketmq.PushConsumer
//var _consumer2 rocketmq.PushConsumer

func init() {

	_producer = newProducer()
}

func newProducer() rocketmq.Producer {
	p, _ := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{EndPoint})),
		producer.WithRetry(2),
	)
	err := p.Start()
	if err != nil {
		util.Panic(err)
	}

	return p
}

//func newConsumerRoom() rocketmq.PushConsumer {
//	c, _ := rocketmq.NewPushConsumer(
//		consumer.WithGroupName(RoomGroup),
//		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{EndPoint})),
//	)
//	err := c.Subscribe(RoomTopic, consumer.MessageSelector{},
//		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
//
//			fmt.Println(util.CurrentSecond(), "Consumer 消费开始 ", RoomTopic, msgs)
//
//			for i := range msgs {
//
//				p := &protocol.Packet{}
//				json.Unmarshal(msgs[i].Body, p)
//
//				if p.Header.Target == protocol.TargetRoom {
//					b, e := service.GetRoom(ctx, p.Header.To)
//
//					if e == nil {
//						for _, v := range b {
//
//							if p.Header.From.UserId == strings.Split(v, "/")[0] {
//								continue
//							}
//
//							broker, _ := service.GetUserDeviceBroker(ctx, v)
//
//							m := protocol.PacketMessage{
//								UserKey: v,
//								Packet:  *p,
//							}
//
//							SendSync2One(broker, &m)
//
//						}
//					}
//
//				}
//
//			}
//			return consumer.ConsumeSuccess, nil
//		})
//
//	if err != nil {
//		util.Panic(err)
//	}
//	// Note: start after subscribe
//	err = c.Start()
//	if err != nil {
//		util.Panic(err)
//	}
//
//	return c
//}
//
//func newConsumerOne() rocketmq.PushConsumer {
//
//	c, _ := rocketmq.NewPushConsumer(
//		consumer.WithGroupName(OneGroup+MyName),
//		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{EndPoint})),
//	)
//	err := c.Subscribe(OneTopic+MyName, consumer.MessageSelector{},
//		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
//
//			fmt.Println(util.CurrentSecond(), "Consumer 消费开始 ", OneTopic+MyName, msgs)
//
//			for i := range msgs {
//				p := &protocol.PacketMessage{}
//				json.Unmarshal(msgs[i].Body, p)
//
//				c, e := service.GetUserContext(p.UserKey)
//
//				if e {
//					serializer.SingleJsonSerializer().Write(c, &p.Packet)
//
//				}
//			}
//			return consumer.ConsumeSuccess, nil
//
//		})
//	if err != nil {
//		util.Panic(err)
//	}
//	// Note: start after subscribe
//	err = c.Start()
//	if err != nil {
//		util.Panic(err)
//	}
//
//	return c
//}

func createTopic(topicName string) {
	endPoint := []string{EndPoint}
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

func sendSync(topic string, message []byte) {
	msg := &primitive.Message{
		Topic: topic,
		Body:  message,
	}
	res, err := _producer.SendSync(context.Background(), msg)

	if err != nil {
		fmt.Printf("send message error: %s\n", err)
	} else {
		fmt.Printf("----------- send message success: result=%s\n", res.String())
	}
}

func SendSync2One(broker string, p *protocol.PacketMessage) {
	msg, _ := json.Marshal(p)
	sendSync(OneTopic+broker2name(broker), msg)
}

func SendSync2Room(packet *protocol.Packet) {
	msg, _ := json.Marshal(packet)
	sendSync(RoomTopic, msg)
}

func broker2name(broker string) string {
	broker = strings.ReplaceAll(broker, ".", "_")
	broker = strings.ReplaceAll(broker, ":", "_")
	return broker
}
