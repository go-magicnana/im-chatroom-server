package mq

import (
	"github.com/segmentio/kafka-go"
	"sync"
)

var once sync.Once

var deliver *kafka.Writer
var deliverOfBroker *kafka.Writer

func DeliverOfAll() *kafka.Writer {
	once.Do(func() {
		deliver = newKafkaWriter("xx", "imchatroom_deliver")
	})
	return deliver
}

func DeliverOfBroker(broker string) *kafka.Writer {
	once.Do(func() {
		deliverOfBroker = newKafkaWriter("xx", "imchatroom_"+broker)
	})
	return deliverOfBroker
}

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}
