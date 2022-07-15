package handler

import (
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"golang.org/x/net/context"
	"im-chatroom-broker/redis"
	"im-chatroom-broker/util"
	"time"
)

const (
	BrokerAlive    string = "imchatroom:brokeralive:"
	BrokerInstance string = "imchatroom:brokerinstance"
)

//func SetBrokerInfo(ctx context.Context, broker, userId string) {
//	redis := redis.Singleton()
//	redis.SAdd(ctx, BrokerInfo+broker, userId)
//}
//
//func DelBrokerInfo(ctx context.Context, broker, userId string) {
//	redis := redis.Singleton()
//	redis.SRem(ctx, BrokerInfo+broker, userId)
//}

func SetBrokerInstance(ctx context.Context, broker string) {
	redis := redis.Singleton()
	redis.SAdd(ctx, BrokerInstance, broker)
}

func DelBrokerInstance(ctx context.Context,broker string){
	redis := redis.Singleton()
	redis.SRem(ctx,BrokerInstance,broker)
}

func SetBrokerAlive(ctx context.Context, broker string) {
	redis := redis.Singleton()
	redis.Set(ctx, BrokerAlive+broker, util.CurrentSecond(), time.Second*70)
}

func GetBrokerAlive(ctx context.Context, broker string) string {
	redis := redis.Singleton()
	cmd := redis.Get(ctx, BrokerAlive+broker)
	if cmd == nil {
		return ""
	}

	return cmd.Val()
}

func BrokerAliveTask(ctx context.Context, broker string) {

	c := cron.New()

																	//0/5 * * * * ? 	每5秒钟1次
	c.AddFunc("*/1 * * * *", func() {			//1分钟1次
		SetBrokerAlive(ctx, broker)
	})

	c.AddFunc("0/30 * * * * ?",func(){
		ProbeBroker(ctx)
	})

	c.Start()
	fmt.Println("BrokerTask Running")

}

func ProbeBroker(ctx context.Context) {
	redis := redis.Singleton()
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
			DelBrokerInstance(ctx,broker)
		}
	}

}

