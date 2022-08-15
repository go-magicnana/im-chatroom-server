package service

import (
	"golang.org/x/net/context"
	"im-chatroom-broker/redis"
)

const (
	BrokerInstance string = "imchatroom:broker.instance"
	BrokerRooms    string = "imchatroom:broker.rooms:"
)

func SetBrokerInstance(broker string) int64 {
	return redis.Rdb.SAdd(context.Background(), BrokerInstance, broker).Val()
}

func GetBrokerInstances() []string {
	redis := redis.Rdb
	return redis.SMembers(context.Background(), BrokerInstance).Val()
}

func SetBrokerRoom(broker, roomId string) int64 {
	return redis.Rdb.SAdd(context.Background(), BrokerRooms+broker, roomId).Val()
}

func GetBrokerRoom(broker string) []string {
	return redis.Rdb.SMembers(context.Background(), BrokerRooms+broker).Val()
}

func GetBrokerRoomExist(broker,roomId string) bool {
	return redis.Rdb.SIsMember(context.Background(),BrokerRooms+broker,roomId).Val()
}

func RemBrokerRoom(broker,roomId string) int64 {
	return redis.Rdb.SRem(context.Background(), BrokerRooms+broker,roomId).Val()
}
