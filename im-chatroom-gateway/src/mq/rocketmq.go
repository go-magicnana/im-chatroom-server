package mq

import (
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"golang.org/x/net/context"
	"im-chatroom-gateway/config"
	"im-chatroom-gateway/protocol"
	"im-chatroom-gateway/util"
	"strings"
)

const (
	RoomGroup = "imchatroom_room_group"
	RoomTopic = "imchatroom_room_topic"

	OneGroup = "imchatroom_one_group_"
	OneTopic = "imchatroom_one_topic_"
)

var _producer rocketmq.Producer

//var _consumer1 rocketmq.PushConsumer
//var _consumer2 rocketmq.PushConsumer

func init() {

	_producer = newProducer()
}

func newProducer() rocketmq.Producer {
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
