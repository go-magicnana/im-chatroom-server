package service

import (
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"golang.org/x/net/context"
	context2 "im-chatroom-broker/context"
	"im-chatroom-broker/redis"
	"im-chatroom-broker/util"
	"time"
)

const (
	BrokerAlive    string = "imchatroom:brokeralive:"
	BrokerCapacity string = "imchatroom:brokercapacity:"
	BrokerInstance string = "imchatroom:brokerinstance"
)

func SetBrokerCapacity(ctx context.Context, broker, userKey string) {
	redis := redis.Rdb
	redis.SAdd(ctx, BrokerCapacity+broker, userKey)
}

func DelBrokerCapacity(ctx context.Context, broker, userKey string) {
	redis := redis.Rdb
	redis.SRem(ctx, BrokerCapacity+broker, userKey)
}

func DelBrokerCapacityAll(ctx context.Context, broker string) {
	redis := redis.Rdb
	redis.Del(ctx, BrokerCapacity+broker)
}

func GetBrokerCapacityAll(ctx context.Context, broker string) []string {
	redis := redis.Rdb
	cmd := redis.SMembers(ctx, BrokerCapacity+broker)
	if cmd != nil {
		clients := cmd.Val()
		if clients != nil && len(clients) > 0 {
			return clients
		}
	}
	return nil
}

func SetBrokerInstance(ctx context.Context, broker string) {
	redis := redis.Rdb
	redis.SAdd(ctx, BrokerInstance, broker)
}

func DelBrokerInstance(ctx context.Context, broker string) {
	redis := redis.Rdb
	redis.SRem(ctx, BrokerInstance, broker)
}

func SetBrokerAlive(ctx context.Context, broker string) {
	redis := redis.Rdb
	redis.Set(ctx, BrokerAlive+broker, util.CurrentSecond(), time.Second*70)
}

func GetBrokerAlive(ctx context.Context, broker string) string {
	redis := redis.Rdb
	cmd := redis.Get(ctx, BrokerAlive+broker)
	if cmd == nil {
		return ""
	}

	return cmd.Val()
}

func AliveTask(ctx context.Context, broker string) {

	c := cron.New()

	//0/5 * * * * ? 	每5秒钟1次
	c.AddFunc("*/1 * * * *", func() { //1分钟1次
		SetBrokerAlive(ctx, broker)
	})

	c.AddFunc("0/30 * * * * ?", func() {
		ProbeBroker(ctx)
	})

	c.AddFunc("*/1 * * * *", func() {
		ProbeConn(ctx)
	})

	c.Start()
	fmt.Println("BrokerTask Running")

}

func ProbeBroker(ctx context.Context) {
	redis := redis.Rdb
	cmd := redis.SMembers(ctx, BrokerInstance)

	if cmd == nil {
		util.Panic(errors.New("无法启动ProbeBroker任务"))
	}

	list, e := cmd.Result()

	if e != nil {
		util.Panic(e)
	}

	for _, broker := range list {
		ret := GetBrokerAlive(ctx, broker)
		if util.IsEmpty(ret) {
			DelBrokerInstance(ctx, broker)
			clients := GetBrokerCapacityAll(ctx, broker)
			DelBrokerCapacityAll(ctx, broker)

			if clients != nil && len(clients) > 0 {
				for _, v := range clients {
					ud, e := GetUserDevice(ctx, v)

					if e == nil && ud != nil {
						DelRoomUser(ctx, ud.RoomId, v)
					}
				}
			}

		}
	}
}

func ProbeConn(ctx context.Context) {
	RangeUserContextAll(func(key, value any) bool {

		user, e := GetUserDevice(ctx, key.(string))

		if e != nil || user == nil {
			c, f := DelUserContext(key.(string))

			if f && c != nil {
				Close(ctx, c)
			}

		}

		return true
	})
}

func Close(ctx context.Context, c *context2.Context) {
	fmt.Println(util.CurrentSecond(), "Read 关闭线程 关闭连接")

	DelUserInfo(ctx, c.UserKey())

	DelUserDevice(ctx, c.UserKey())

	DelRoomUser(ctx, c.RoomId(), c.UserKey())

	DelUserContext(c.UserKey())

	DelBrokerCapacity(ctx, c.Broker(), c.UserKey())

	c.Close()

	c = nil

}
