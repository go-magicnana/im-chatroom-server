package mq

import (
	"encoding/json"
	"fmt"
	"gitee.com/zhucheer/orange/queue"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestMq(t *testing.T) {
	// 注册生产者 填入broker节点,group名称,重试次数信息
	mqProducerClient := queue.RegisterRocketProducerMust([]string{"192.168.3.242:9876"}, "test", 1)

	// 注册消费者 填入broker节点,group名称信息
	mqConsumerClient := queue.RegisterRocketConsumerMust([]string{"192.168.3.242:9876"}, "test")

	go func() {
		for i := 0; i < 10; i++ {
			// 向队列发送一条消息 填入消息队列topic和消息体信息
			ret, _ := mqProducerClient.SendMsg("topicTest", "Hello mq~~")
			fmt.Println("========producer push one message====", ret.MsgId)

			time.Sleep(time.Second)
		}

	}()

	// 执行消费者监听 填入消息队列topic
	mqConsumerClient.ListenReceiveMsgDo("topicTest", func(mqMsg queue.MqMsg) {
		// 收到一条消息
		fmt.Println("receive====>", mqMsg.MsgId, mqMsg.BodyString())

	})

	time.Sleep(20 * time.Second)

}

func TestRocket(t *testing.T) {

	val, _ := json.Marshal("jsonMessage")
	message := primitive.NewMessage("imchatroom_deliver", val)
	result, err := Deliver().SendSync(context.Background(), message)
	if err != nil {
		return
	}
	fmt.Println(result)

	c := NewRocketMqConsumer()
	c.Subscribe("imchatroom_deliver", consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			fmt.Printf("subscribe callback : %v \n", msgs[i])
		}
		return consumer.ConsumeSuccess, nil
	})
	c.Start()
}
