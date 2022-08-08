package service

import (
	"golang.org/x/net/context"
	"im-chatroom-broker/redis"
)

const (
	BrokerInstance string = "imchatroom:broker.instance"
)

func SetBrokerInstance(broker string) int64 {
	return redis.Rdb.SAdd(context.Background(), BrokerInstance, broker).Val()
}

func GetBrokerInstances() []string {
	redis := redis.Rdb
	return redis.SMembers(context.Background(), BrokerInstance).Val()
}
