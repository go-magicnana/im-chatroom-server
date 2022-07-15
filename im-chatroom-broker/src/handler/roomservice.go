package handler

import (
	"golang.org/x/net/context"
	"im-chatroom-broker/redis"
)

const (
	RoomInfo string = "imchatroom:roominfo:"
)

func SetRoomUser(ctx context.Context, roomId string, userKey string) {
	redis := redis.Singleton()
	redis.SAdd(ctx, RoomInfo+roomId, userKey)
}

//func GetRoom(ctx context.Context, roomId string) (*[...]string, error){
//	redis := redis.Singleton()
//	cmd := redis.SCard(ctx, RoomInfo+roomId)
//	m, e := cmd.Result()
//	if e != nil {
//		return nil,e
//	}
//	fmt.Println(m)
//	return m,nil
//}

func DelRoomUser(ctx context.Context, roomId string, userKey string) {
	redis := redis.Singleton()
	redis.SRem(ctx, RoomInfo+roomId, userKey)
}

func ClearRoom(ctx context.Context, roomId string) {
	redis := redis.Singleton()
	redis.Del(ctx, RoomInfo+roomId)
}
