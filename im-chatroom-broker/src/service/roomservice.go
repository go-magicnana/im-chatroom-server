package service

import (
	"golang.org/x/net/context"
	"im-chatroom-broker/redis"
)

const (
	RoomInfo string = "imchatroom:roominfo:"
)

func SetRoomUser(ctx context.Context, roomId string, userKey string) {
	redis.Rdb.SAdd(ctx, RoomInfo+roomId, userKey)
}

func GetRoom(ctx context.Context, roomId string) ([]string, error) {
	cmd := redis.Rdb.SMembers(ctx, RoomInfo+roomId)
	m, e := cmd.Result()
	if e != nil {
		return nil, e
	}
	return m, nil
}

func DelRoomUser(ctx context.Context, roomId string, userKey string) {
	redis.Rdb.SRem(ctx, RoomInfo+roomId, userKey)

	userKeys, _ := GetRoom(ctx, roomId)
	if userKeys == nil || len(userKeys) == 0 {
		DelRoom(ctx, roomId)
	}
}

func DelRoom(ctx context.Context, roomId string) {
	redis.Rdb.Del(ctx, RoomInfo+roomId)
}
