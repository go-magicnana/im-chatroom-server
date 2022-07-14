package handler

import (
	"golang.org/x/net/context"
	"im-chatroom-broker/redis"
)

const (
	// BrokerInfo SET
	BrokerInfo string = "imchatroom:brokerinfo:"


	BrokerInstance string = "imchatroom:brokerinstance"
)

func SetUserBroker(ctx context.Context, broker, userId string) {
	redis := redis.Singleton()
	redis.SAdd(ctx, BrokerInfo+broker, userId)
}

func SetBrokerInfo(ctx context.Context,broker string){
	redis := redis.Singleton()
	redis.SAdd(ctx,BrokerInstance,broker)
}
