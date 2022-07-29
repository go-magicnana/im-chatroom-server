package service

import (
	"github.com/robfig/cron/v3"
	"golang.org/x/net/context"
	context2 "im-chatroom-broker/context"
	"im-chatroom-broker/redis"
	"im-chatroom-broker/zaplog"
)

const (
	BrokerCapacity string = "imchatroom:broker.capacity:"
	BrokerInstance string = "imchatroom:broker.instance"
)

func SetBrokerCapacity(ctx context.Context, broker, clientName string) {
	redis := redis.Rdb
	redis.SAdd(ctx, BrokerCapacity+broker, clientName)
}

func DelBrokerCapacity(ctx context.Context, broker, clientName string) {
	redis := redis.Rdb
	redis.SRem(ctx, BrokerCapacity+broker, clientName)
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

//func SetBrokerAlive(ctx context.Context, broker string) {
//	redis := redis.Rdb
//	redis.Set(ctx, BrokerAlive+broker, util.CurrentSecond(), time.Second*70)
//}
//
//func GetBrokerAlive(ctx context.Context, broker string) string {
//	redis := redis.Rdb
//	cmd := redis.Get(ctx, BrokerAlive+broker)
//	if cmd == nil {
//		return ""
//	}
//
//	return cmd.Val()
//}
//
func AliveTask(ctx context.Context, broker string) {

	c := cron.New()

	////0/5 * * * * ? 	每5秒钟1次
	//c.AddFunc("*/1 * * * *", func() { //1分钟1次
	//	SetBrokerAlive(ctx, broker)
	//	zaplog.Logger.Debugf("Task SetBrokerAlive %s", broker)
	//
	//})

	//c.AddFunc("@every 1s", func() {
	//	//ProbeBroker(ctx)
	//	SetBrokerInstance(ctx, broker)
	//	zaplog.Logger.Debugf("Task SetBrokerInstance %s", broker)
	//})

	//c.AddFunc("@every 1m", func() {
	//	ProbeConn(ctx)
	//	zaplog.Logger.Debugf("Task ProbeConns %s", broker)
	//})
	//
	//c.AddFunc("@every 1m", func() {
	//	ProbeRoom(ctx)
	//	zaplog.Logger.Debugf("Task ProbeRoom %s", broker)
	//})

	c.Start()
	zaplog.Logger.Infof("Task Running %s", broker)

}

//func ProbeRoom(ctx context.Context) {
//	roomList := GetRoomInstance(ctx)
//
//	if roomList != nil && len(roomList) > 0 {
//		for _, v := range roomList {
//			members, err := GetRoomMembers(ctx, v)
//
//			if (err == nil || err == redis.Nil) && len(members) < 10 && len(members) > 0 {
//				for _, vv := range members {
//					userDevice, e := GetUserDevice(ctx, vv)
//
//					if (e == nil || e == redis.Nil) && userDevice == nil {
//						DelRoomUser(ctx, v, vv)
//					}
//				}
//			}
//		}
//	}
//
//}

//func ProbeBroker(ctx context.Context) {
//	redis := redis.Rdb
//	cmd := redis.SMembers(ctx, BrokerInstance)
//
//	if cmd == nil {
//		util.Panic(errors.New("无法启动ProbeBroker任务"))
//	}
//
//	list, e := cmd.Result()
//
//	if e != nil {
//		util.Panic(e)
//	}
//
//	for _, broker := range list {
//		ret := GetBrokerAlive(ctx, broker)
//		if util.IsEmpty(ret) {
//			clients := GetBrokerCapacityAll(ctx, broker)
//
//			if clients != nil && len(clients) > 0 {
//				for _, v := range clients {
//					ud, e := GetUserDevice(ctx, v)
//
//					if e == nil && ud != nil {
//						DelRoomUser(ctx, ud.RoomId, v)
//					}
//				}
//			}
//
//			DelBrokerInstance(ctx, broker)
//			DelBrokerCapacityAll(ctx, broker)
//		}
//	}
//}

//func ProbeConn(ctx context.Context) {
//	RangeUserContextAll(func(key, value any) bool {
//
//		user, e := GetUserDevice(ctx, key.(string))
//
//		if e != nil || user == nil {
//			c, f := DelUserContext(key.(string))
//
//			if f && c != nil {
//				Close(ctx, c)
//			}
//
//		}
//
//		return true
//	})
//}

func Close(ctx context.Context, c *context2.Context) {

	//DelUserInfo(ctx, c.UserId())

	//DelUserDevice(ctx, c.ClientName())

	DelUserClient(ctx, c.UserId(), c.ClientName())

	//DelRoomUser(ctx, c.RoomId(), c.ClientName())

	DelUserContext(c.ClientName())

	DelRoomClients(c.RoomId(), c.ClientName())

	DelBrokerCapacity(ctx, c.Broker(), c.ClientName())

	zaplog.Logger.Infof("CloseByClient %s", c.Conn().RemoteAddr())

	c.Close()

	c = nil
}
