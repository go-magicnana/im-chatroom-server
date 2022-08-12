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
	"im-chatroom-broker/deliver"
	"im-chatroom-broker/protocol"
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

//var _consumer1 rocketmq.PushConsumer
var _consumer2 rocketmq.PushConsumer

func init() {

	if util.IsNotEmpty(config.OP.Ip) {
		MyName = broker2name(config.OP.Ip + ":" + config.OP.Port)

	} else {
		ip := util.GetBrokerIp()
		MyName = broker2name(ip + ":" + config.OP.Port)

	}
	//
	//PushRoomChannel = make(chan *protocol.PacketMessage, 500)
	//Push2OneChannel = make(chan *protocol.PacketMessage, 1500)
	_producer = newProducer()
	//_consumer1 = newConsumerRoom()
	_consumer2 = newConsumerOne()
	//
	//c, _ := context.WithCancel(context.Background())
	//
	//go SendSync2Room(c)
	//go SendSync2One(c)

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

//func newConsumerRoom() rocketmq.PushConsumer {
//
//	c, _ := rocketmq.NewPushConsumer(
//		consumer.WithGroupName(RoomGroup),
//		consumer.WithConsumerModel(consumer.BroadCasting),
//		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{config.OP.RocketMQ.Address})),
//	)
//	err := c.Subscribe(RoomTopic, consumer.MessageSelector{},
//		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
//
//			for i := range msgs {
//
//				p := &protocol.PacketMessage{}
//				json.Unmarshal(msgs[i].Body, p)
//
//				if p.Broker == thread.BrokerAddress {
//					continue
//				}
//
//				zaplog.Logger.Debugf("CsumRoom %s %s C:%d T:%d F:%d %v %v", RoomTopic, p.Packet.Header.MessageId, p.Packet.Header.Command, p.Packet.Header.Type, p.Packet.Header.Flow, p.Packet.Body, msgs[i].MsgId)
//
//				if p.Packet.Header.Target == protocol.TargetRoom {
//
//					roomClients := thread.GetRoomChannels(p.Packet.Header.To)
//
//					for _, clientName := range roomClients {
//						cc := thread.GetChannel(clientName.(string))
//
//						if cc != nil && cc.Channel != nil {
//							cc.Channel <- protocol.NewResponse(p.Packet)
//						}
//					}
//
//				}
//
//			}
//			zaplog.Logger.Debugf("---------------------------------------------------------")
//
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
//	zaplog.Logger.Infof("Init RocketMQ Consumer-Room %s", config.OP.RocketMQ.Address)
//
//	return c
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

				deliver.Deliver2Worker(p.Broker, p.Packet)

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

func broker2name(broker string) string {
	broker = strings.ReplaceAll(broker, ".", "_")
	broker = strings.ReplaceAll(broker, ":", "_")
	return broker
}

func Deliver2MQ(localBroker string, packet *protocol.Packet) {

	if packet.Header.Target == protocol.TargetRoom {
		brokers := service.GetBrokerInstances()

		for _, broker := range brokers {
			if localBroker == broker {
				continue
			}

			size := service.CardRoomClients(broker, packet.Header.To)
			if size <= 0 {
				continue
			}

			p := protocol.PacketMessage{
				Broker: broker,
				Packet: packet,
			}

			topic := OneTopic + broker2name(broker)
			msg, _ := json.Marshal(p)
			sendSync(topic, msg)

		}
	}
}
