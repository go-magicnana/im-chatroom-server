package service

import (
	"golang.org/x/net/context"
	"im-chatroom-gateway/redis"
	"im-chatroom-gateway/util"
)

const (
	LockBrokerHeartbeat string = "imchatroom:lock.heartbeat:"
)

func Lock(ctx context.Context, broker string) bool {
	s := redis.Rdb.SetNX(ctx, LockBrokerHeartbeat+broker, util.CurrentSecond(),-1).Val()
	return s
}

func Unlock(ctx context.Context, broker string) {
	redis.Rdb.Del(ctx, LockBrokerHeartbeat+broker)
}
