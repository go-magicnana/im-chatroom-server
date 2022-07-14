package handler

import (
	"golang.org/x/net/context"
	"im-chatroom-broker/redis"
)

const (
	BrokerInfo string = "imchatroom:brokerinfo:"
	BrokerInstance string = "imchatroom:brokerinstance"
)

func SetBrokerInfo(ctx context.Context, broker, userId string) {
	redis := redis.Singleton()
	redis.SAdd(ctx, BrokerInfo+broker, userId)
}

func DelBrokerInfo(ctx context.Context,broker,userId string){
	redis := redis.Singleton()
	redis.SRem(ctx,BrokerInfo+broker,userId)
}

func SetBrokerInstance(ctx context.Context,broker string){
	redis := redis.Singleton()
	redis.SAdd(ctx,BrokerInstance,broker)
}
