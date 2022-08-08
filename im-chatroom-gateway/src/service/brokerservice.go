package service

import (
	"golang.org/x/net/context"
	"im-chatroom-gateway/redis"
)

const (
	BrokerInstance string = "imchatroom:broker.instance"
)

func GetBrokerInstances(ctx context.Context) []string {
	redis := redis.Rdb
	return redis.SMembers(ctx, BrokerInstance).Val()
}

func DelBrokerInstance(ctx context.Context, broker string) {
	redis := redis.Rdb
	redis.SRem(ctx, BrokerInstance, broker)
}

