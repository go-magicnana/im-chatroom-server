package service

import (
	"golang.org/x/net/context"
	"im-chatroom-broker/redis"
	"strconv"
)

const (
	// hash
	RoomInfo string = "imchatroom:room.info:"
	// set
	RoomMembers string = "imchatroom:room.members:"
	RoomBlacks  string = "imchatroom:room.blacks:"
)

func SetRoomUser(ctx context.Context, roomId string, userKey string) {
	redis.Rdb.SAdd(ctx, RoomMembers+roomId, userKey)
}

func GetRoomMembers(ctx context.Context, roomId string) ([]string, error) {
	cmd := redis.Rdb.SMembers(ctx, RoomMembers+roomId)
	m, e := cmd.Result()
	if e != nil {
		return nil, e
	}
	return m, nil
}

func DelRoomUser(ctx context.Context, roomId string, userKey string) {
	redis.Rdb.SRem(ctx, RoomMembers+roomId, userKey)
}

func GetRoomBlocked(ctx context.Context, roomId string) int {
	cmd := redis.Rdb.HGet(ctx, RoomInfo+roomId, "blocked")
	result, err := cmd.Result()
	if err != nil {
		return 0
	}
	atoi, _ := strconv.Atoi(result)
	return atoi
}

func GetRoomMemberBlocked(ctx context.Context, roomId string, userId string) int {
	cmd := redis.Rdb.SIsMember(ctx, RoomBlacks+roomId, userId)
	m, e := cmd.Result()
	if e != nil || m == false {
		return 0
	}
	return 1
}
