package server

import (
	"context"
	"fmt"
	"im-chatroom-broker/redis"
	"im-chatroom-broker/util"
	"strconv"
)

const (
	BrokerHash string = "dudu:broker:broker_hash"
)

func SetBroker(ctx context.Context,brokerId int,address string)  {
	cmd1 := redis.Singleton().HSet(ctx, BrokerHash,strconv.Itoa(brokerId),address)
	if cmd1.Err()!=nil {
		util.Panic(cmd1.Err())
	}
	fmt.Println(cmd1)
}

func GetBrokers(ctx context.Context) map[string]string{
	ret := redis.Singleton().HGetAll(ctx, BrokerHash)

	if ret.Err()!=nil {
		util.Panic(ret.Err())
	}

	if ret!=nil{
		return ret.Val()
	}
	return nil
}


