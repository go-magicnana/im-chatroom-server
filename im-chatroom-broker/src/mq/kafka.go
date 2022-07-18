package mq

//
//import (
//	"encoding/json"
//	"github.com/segmentio/kafka-go"
//	"golang.org/x/net/context"
//	context2 "im-chatroom-broker/context"
//	"im-chatroom-broker/protocol"
//	"sync"
//)
//
//var once sync.Once
//
//var client *kafka.Writer
//
//func deliver() *kafka.Writer {
//	once.Do(func() {
//		client = newKafkaWriter("xx")
//	})
//	return client
//}
//
//func newKafkaWriter(kafkaURL string) *kafka.Writer {
//	return &kafka.Writer{
//		Addr:     kafka.TCP(kafkaURL),
//		Balancer: &kafka.LeastBytes{},
//	}
//}
//
//func DeliverMessageToRoom(ctx context.Context, c *context2.Context, packet *protocol.Packet) {
//
//	val, e := json.Marshal(packet)
//
//	if e != nil {
//		return
//	}
//
//	if len(val) == 0 {
//		return
//	}
//
//	msg := kafka.Message{
//		Key:   []byte(packet.Header.To),
//		Value: val,
//		Topic: "imchatroom_deliver",
//	}
//
//	e = deliver().WriteMessages(ctx, msg)
//	if e != nil {
//		return
//	}
//}
//
//func DeliverMessageToUser(ctx context.Context, c *context2.Context, packet *protocol.Packet) {
//
//	val, e := json.Marshal(packet)
//
//	if e != nil {
//		return
//	}
//
//	if len(val) == 0 {
//		return
//	}
//
//	msg := kafka.Message{
//		Key:   []byte(packet.Header.To),
//		Value: val,
//		Topic: "imchatroom_" + c.Broker(),
//	}
//
//	e = deliver().WriteMessages(ctx, msg)
//	if e != nil {
//		return
//	}
//}
