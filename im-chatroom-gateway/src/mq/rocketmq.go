package mq

import (
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"golang.org/x/net/context"
	"im-chatroom-gateway/src/protocol"
	"im-chatroom-gateway/src/util"
	"strings"
	"sync"
)

var once sync.Once

var _mq Deliver

type Deliver struct {
	Producer rocketmq.Producer
	Consumer rocketmq.PushConsumer
}

var (
	mqProducer rocketmq.Producer
	mqConsumer rocketmq.PushConsumer
)

func OneDeliver() Deliver {
	once.Do(func() {

		c, _ := rocketmq.NewPushConsumer(
			//consumer.WithGroupName("testGroup"),
			consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"192.168.3.242:9876"})),
		)

		p, _ := rocketmq.NewProducer(
			producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"192.168.3.242:9876"})),
			producer.WithRetry(2),
		)

		p.Start()

		_mq = Deliver{
			Producer: p,
			Consumer: c,
		}
	})
	return _mq
}

func (d Deliver) Sync(topic string, body []byte) {

	msg := &primitive.Message{
		Topic: topic,
		Body:  body,
	}
	res, err := d.Producer.SendSync(context.Background(), msg)

	fmt.Println(util.CurrentSecond(), "Producer 发送完毕", topic, res, err)

}

func (d Deliver) ProduceRoom(packet *protocol.Packet) {
	msg, _ := json.Marshal(packet)
	d.Sync("imchatroom_push_room", msg)
}

func (d Deliver) ProduceOne(broker string, packet *protocol.PacketMessage) {
	msg, _ := json.Marshal(packet)

	broker = strings.ReplaceAll(broker, ".", "_")
	broker = strings.ReplaceAll(broker, ":", "_")

	d.Sync("imchatroom_push_one_"+broker, msg)
	fmt.Println(util.CurrentSecond())
}

//func (d Deliver) consume(topic string, f func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)) {
//	err := d.Consumer.Subscribe(topic, consumer.MessageSelector{}, f)
//	if err != nil {
//		util.Panic(err)
//	}
//
//	d.Consumer.Start()
//}

//func (d Deliver) ConsumeRoom() {
//	d.consume("imchatroom_push_room", func(c context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
//
//		fmt.Println(util.CurrentSecond(), "Consumer 消费开始 imchatroom_push_room", msgs)
//
//		for i := range msgs {
//
//			p := &protocol.Packet{}
//			json.Unmarshal(msgs[i].Body, p)
//
//			if p.Header.Target == protocol.TargetRoom {
//				b, e := service.GetRoom(c, p.Header.To)
//
//				if e == nil {
//					for _, v := range b {
//
//						if p.Header.From.UserId == strings.Split(v, "/")[0] {
//							continue
//						}
//
//						broker, _ := service.GetUserDeviceBroker(c, v)
//
//						m := protocol.PacketMessage{
//							UserKey: v,
//							Packet:  *p,
//						}
//						d.ProduceOne(broker, &m)
//					}
//				}
//
//			}
//
//		}
//		return consumer.ConsumeSuccess, nil
//	})
//}
//
//func (d Deliver) ConsumeMine(broker string) {
//
//	broker = strings.ReplaceAll(broker, ".", "_")
//	broker = strings.ReplaceAll(broker, ":", "_")
//
//	d.consume("imchatroom_push_one_"+broker, func(c context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
//
//		fmt.Println(util.CurrentSecond(), "Consumer 消费开始 imchatroom_push_one_"+broker, msgs)
//
//		for i := range msgs {
//			p := &protocol.PacketMessage{}
//			json.Unmarshal(msgs[i].Body, p)
//
//			c, e := service.GetUserContext(p.UserKey)
//
//			if e {
//				//c.Push(&p.Packet)
//				serializer.SingleJsonSerializer().Write(c, &p.Packet)
//
//			}
//		}
//		return consumer.ConsumeSuccess, nil
//	})
//}
