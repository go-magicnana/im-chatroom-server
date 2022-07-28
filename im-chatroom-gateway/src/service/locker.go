package service

import (
	"golang.org/x/net/context"
	"im-chatroom-gateway/redis"
	"im-chatroom-gateway/util"
	"im-chatroom-gateway/zaplog"
)

const (
	LockBrokerHeartbeat string = "imchatroom:lock.heartbeat:"
)

func Lock(ctx context.Context, broker string) bool {
	s := redis.Rdb.SetNX(ctx, LockBrokerHeartbeat+broker, util.CurrentSecond(),-1).Val()
	zaplog.Logger.Debugf("Heartbeat %s Lock %v", broker,s)
	return s
}

func Unlock(ctx context.Context, broker string) {
	s := redis.Rdb.Del(ctx, LockBrokerHeartbeat+broker).Val()
	zaplog.Logger.Debugf("Heartbeat %s Unlock %v", broker,s)

}
