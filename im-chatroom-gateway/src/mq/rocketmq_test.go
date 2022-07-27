package mq

import (
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/hashicorp/go-uuid"
	"golang.org/x/net/context"
	"im-chatroom-gateway/protocol"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

func SendSyncMessage(topic, message string) {
	p, _ := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"192.168.3.242:9876"})),
		producer.WithRetry(2),
	)
	err := p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}

	for i := 0; i < 10; i++ {
		msg := &primitive.Message{
			Topic: topic,
			Body:  []byte("Hello RocketMQ Go Client! " + strconv.Itoa(i)),
		}
		res, err := p.SendSync(context.Background(), msg)

		if err != nil {
			fmt.Printf("send message error: %s\n", err)
		} else {
			fmt.Printf("----------- send message success: result=%s\n", res.String())
		}
	}
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}
}

func SubscribeMessage(topic, group string) {
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName(group),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"192.168.3.242:9876"})),
	)
	err := c.Subscribe(topic, consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			fmt.Printf("----------- subscribe callback: %v \n", msgs[i])
		}

		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	// Note: start after subscribe
	err = c.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	time.Sleep(time.Hour)
	err = c.Shutdown()
	if err != nil {
		fmt.Printf("shutdown Consumer error: %s", err.Error())
	}
}

func TestRocket(t *testing.T) {

	var wg sync.WaitGroup

	wg.Add(1)

	topic1 := "test1Topic"
	group1 := "test1Group"

	topic2 := "test2Topic"
	group2 := "test2Group"

	go SubscribeMessage(topic1, group1)
	go SubscribeMessage(topic2, group2)

	for {
		i := 0
		SendSyncMessage(topic1, "hello world "+strconv.Itoa(i))
		time.Sleep(time.Second * 5)
		SendSyncMessage(topic2, "hello world "+strconv.Itoa(i))

		i++

	}

	wg.Wait()

}

func TestMQ2(t *testing.T) {

	messageId, _ := uuid.GenerateUUID()

	//TypeNoticeBlockUser   = 3103
	//TypeNoticeUnblockUser = 3104
	//TypeNoticeCloseRoom   = 3105
	//TypeNoticeBlockRoom   = 3106
	//TypeNoticeUnblockRoom = 3107

	header := protocol.MessageHeader{
		MessageId: messageId,
		Command:   protocol.CommandNotice,
		Target:    1,
		From: protocol.UserInfo{
			UserId: "110",
			Name:   "wgw",
			Avatar: "http://www.baidu.com",
		},
		To:   "1",
		Flow: protocol.FlowDown,
		Type: 3106,
	}

	var body any

	body = protocol.MessageBodyNoticeBlockRoom{RoomId: "233"}

	packet := protocol.Packet{
		Header: header, Body: body,
	}

	SendSync2Room(&packet)
}
