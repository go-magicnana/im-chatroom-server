package service

import (
	"golang.org/x/net/context"
	"im-chatroom-broker/redis"
)

const (
	UserClientMessage string = "test:receive"
)

func AddUserClientMessage(ctx context.Context, clientName string, num int64) int64 {
	ret := redis.Rdb.HIncrBy(ctx, UserClientMessage, clientName, num).Val()
	return ret
}
