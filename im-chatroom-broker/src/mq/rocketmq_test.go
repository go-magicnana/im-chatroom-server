package mq

import (
	"sync"
	"testing"
)

func TestMq(t *testing.T) {
	//// 注册生产者 填入broker节点,group名称,重试次数信息
	//mqProducerClient := queue.RegisterRocketProducerMust([]string{"192.168.3.242:9876"}, "test", 1)
	//
	//// 注册消费者 填入broker节点,group名称信息
	//mqConsumerClient := queue.RegisterRocketConsumerMust([]string{"192.168.3.242:9876"}, "test")
	//
	//go func() {
	//	for i := 0; i < 10; i++ {
	//		// 向队列发送一条消息 填入消息队列topic和消息体信息
	//		ret, _ := mqProducerClient.SendMsg("topicTest", "Hello mq~~")
	//		fmt.Println("========producer push one message====", ret.MsgId)
	//
	//		time.Sleep(time.Second)
	//	}
	//
	//}()
	//
	//// 执行消费者监听 填入消息队列topic
	//mqConsumerClient.ListenReceiveMsgDo("topicTest", func(mqMsg queue.MqMsg) {
	//	// 收到一条消息
	//	fmt.Println("receive====>", mqMsg.MsgId, mqMsg.BodyString())
	//
	//})
	//
	//time.Sleep(20 * time.Second)

}

func TestRocket(t *testing.T) {

	var wg sync.WaitGroup

	wg.Add(1)

	//OneDeliver().ConsumeRoom()
	//OneDeliver().ConsumeMine("192.168.3.92:33121")

	//p := &protocol.Packet{
	//	Header: protocol.MessageHeader{
	//		Command: 9,
	//	},
	//	Body: "JJJJJ",
	//}
	//OneDeliver().ProduceRoom(p)

	//pac := &protocol.PacketMessage{
	//	Packet: protocol.Packet{
	//		Header: protocol.MessageHeader{
	//			To: "106533",
	//		},
	//		Body: "hahaha",
	//	},
	//	UserKey: "106533/WW",
	//}
	//OneDeliver().ProduceOne("192.168.3.92:33121", pac)
	//			Command: 4,
	//			Target: 2,
	//
	//			To: "1066",
	//			Flow: 1,
	//		},
	//		Body: "hahaha",
	//	},
	//	UserKey: "1066/IOS",
	//}
	//OneDeliver().ProduceOne("192.168.3.87:33121", pac)

	wg.Wait()

}
