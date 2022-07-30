package service

import (
	"github.com/robfig/cron/v3"
	"golang.org/x/net/context"
	"im-chatroom-broker/redis"
	"im-chatroom-broker/zaplog"
)

const (
	//BrokerClients  string = "imchatroom:broker.clients:"
	BrokerInstance string = "imchatroom:broker.instance"
)

//
//func SetBrokerClients(ctx context.Context, broker, clientName string) int64 {
//	return redis.Rdb.SAdd(ctx, BrokerClients+broker, clientName).Val()
//}
//
//func RemBrokerClients(ctx context.Context, broker, clientName string) int64 {
//	return redis.Rdb.SRem(ctx, BrokerClients+broker, clientName).Val()
//}
//
//func DelBrokerClients(ctx context.Context, broker string) int64 {
//	return redis.Rdb.Del(ctx, BrokerClients+broker).Val()
//}
//
//func GetBrokerClients(ctx context.Context, broker string) []string {
//	return redis.Rdb.SMembers(ctx, BrokerClients+broker).Val()
//}

func SetBrokerInstance(ctx context.Context, broker string) int64 {
	return redis.Rdb.SAdd(ctx, BrokerInstance, broker).Val()
}

func RemBrokerInstance(ctx context.Context, broker string) int64 {
	return redis.Rdb.SRem(ctx, BrokerInstance, broker).Val()
}

func AliveTask(ctx context.Context, broker string) {

	defer func() {
		zaplog.Logger.Errorf("TaskRecover %v", recover())
	}()

	c := cron.New()

	////0/5 * * * * ? 	每5秒钟1次
	//c.AddFunc("*/1 * * * *", func() { //1分钟1次
	//	SetBrokerAlive(ctx, broker)
	//	zaplog.Logger.Debugf("Task SetBrokerAlive %s", broker)
	//
	//})

	c.AddFunc("@every 1s", func() {
		//ProbeBroker(ctx)
		SetBrokerInstance(ctx, broker)
		//zaplog.Logger.Debugf("Task SetBrokerInstance %s", broker)
	})

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
