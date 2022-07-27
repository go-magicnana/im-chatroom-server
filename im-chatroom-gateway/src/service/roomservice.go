package service

import (
	"golang.org/x/net/context"
	"im-chatroom-gateway/redis"
)

const (
	// hash
	RoomInfo string = "imchatroom:room.info:"
	// set
	RoomMembers  string = "imchatroom:room.members:"
	RoomBlacks   string = "imchatroom:room.blacks:"
	RoomInstance string = "imchatroom:room.instance"
)

func SetRoomInstance(ctx context.Context, roomId string) int64 {
	return redis.Rdb.SAdd(ctx, RoomInstance, roomId).Val()
}

// block 0正常；1封禁
func SetRoomBlocked(ctx context.Context, roomId string, block int) int64 {
	cmd := redis.Rdb.HSet(ctx, RoomInfo+roomId, "blocked", block)
	result, err := cmd.Result()
	if err != nil {
		return 0
	}
	return result
}

func SetRoomMemberBlocked(ctx context.Context, roomId string, userId string) int {
	cmd := redis.Rdb.SAdd(ctx, RoomBlacks+roomId, userId)
	m, e := cmd.Result()
	if e != nil || m == 0 {
		return 0
	}
	return 1
}

func RemRoomMemberBlocked(ctx context.Context, roomId string, userId string) int {
	cmd := redis.Rdb.SRem(ctx, RoomBlacks+roomId, userId)
	m, e := cmd.Result()
	if e != nil || m == 0 {
		return 0
	}
	return 1
}
