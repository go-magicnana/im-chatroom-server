package service

import (
	"golang.org/x/net/context"
	"im-chatroom-broker/redis"
)

const (
	UserClientMessage string = "test:user.client.message:"
)

func AddUserClientMessage(ctx context.Context, clientName string, msg string) int64 {
	return redis.Rdb.RPush(ctx, UserClientMessage+clientName, msg).Val()
}
