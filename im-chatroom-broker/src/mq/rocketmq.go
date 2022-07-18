package mq

import (
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"golang.org/x/net/context"
	context2 "im-chatroom-broker/context"
	"im-chatroom-broker/protocol"
	"os"
	"strconv"
	"sync"
)

var once sync.Once

var (
	mqProducer rocketmq.Producer
	mqConsumer rocketmq.PushConsumer
)

func Deliver() rocketmq.Producer {
	once.Do(func() {
		mqProducer = newRocketMqProducer()
	})
	return mqProducer
}

func newRocketMqProducer() rocketmq.Producer {
	endPoint := []string{"192.168.3.242:9876"}
	newProducer, _ := rocketmq.NewProducer(
		producer.WithNameServer(endPoint),
		producer.WithRetry(3),
		producer.WithGroupName("ProducerGroupName"),
	)
	err := newProducer.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}
	return newProducer
}

func Consumer() rocketmq.PushConsumer {
	once.Do(func() {
		mqConsumer = NewRocketMqConsumer()
	})
	return mqConsumer
}

func NewRocketMqConsumer() rocketmq.PushConsumer {
	endPoint := []string{"192.168.3.242:9876"}
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer(endPoint),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithGroupName("ProducerGroupName"),
	)
	return c
}

func DeliverMessageToRoom(ctx context.Context, c *context2.Context, packet *protocol.Packet) {

	val, e := json.Marshal(packet)

	if e != nil {
		return
	}

	if len(val) == 0 {
		return
	}

	message := primitive.NewMessage("imchatroom_deliver", val).WithTag(strconv.Itoa(int(packet.Header.Target)))

	result, err := Deliver().SendSync(ctx, message)
	if err != nil {
		return
	}
	fmt.Println(result)
}

func DeliverMessageToUser(ctx context.Context, c *context2.Context, packet *protocol.Packet) {

	val, e := json.Marshal(packet)

	if e != nil {
		return
	}

	if len(val) == 0 {
		return
	}

	msg := primitive.NewMessage("imchatroom_deliver", val).WithTag(strconv.Itoa(int(packet.Header.Target)))

	result, err := Deliver().SendSync(ctx, msg)
	if err != nil {
		return
	}
	fmt.Println(result)
}
