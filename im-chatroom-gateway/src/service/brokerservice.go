package service

import (
	"golang.org/x/net/context"
	"im-chatroom-gateway/redis"
)

const (
	BrokerAlive    string = "imchatroom:broker.alive:"
	BrokerCapacity string = "imchatroom:broker.capacity:"
	BrokerInstance string = "imchatroom:broker.instance"
)

func DelBrokerCapacityAll(ctx context.Context, broker string) {
	redis := redis.Rdb
	redis.Del(ctx, BrokerCapacity+broker)
}

func GetBrokerInstance(ctx context.Context) []string {
	redis := redis.Rdb
	return redis.SMembers(ctx, BrokerInstance).Val()
}

func DelBrokerInstance(ctx context.Context, broker string) {
	redis := redis.Rdb
	redis.SRem(ctx, BrokerInstance, broker)
}

