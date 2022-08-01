package service

import (
	"golang.org/x/net/context"
	"im-chatroom-broker/redis"
	"time"
)

const (
	UserClientMessage string = "test:user.client.message:"
)

func AddUserClientMessage(ctx context.Context, clientName string, msg string) int64 {
	ret := redis.Rdb.RPush(ctx, UserClientMessage+clientName, msg).Val()

	if msg == "99" {
		redis.Rdb.Expire(ctx, UserClientMessage+clientName, time.Minute*3)
	}

	return ret
}
