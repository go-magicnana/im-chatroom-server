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
	RoomInstance  string = "imchatroom:room.instance"
)


//func GetRoomInstance(ctx context.Context) []string{
//	return redis.Rdb.SMembers(ctx,RoomInstance).Val()
//}
//
//func DelRoomInstance(ctx context.Context,roomId string){
//	redis.Rdb.SRem(ctx,RoomInstance,roomId)
//}

//func SetRoomUser(ctx context.Context, roomId string, clientName string) {
//	redis.Rdb.SAdd(ctx, RoomMembers+roomId, clientName)
//}
//
//func GetRoomMembers(ctx context.Context, roomId string) ([]string, error) {
//	cmd := redis.Rdb.SMembers(ctx, RoomMembers+roomId)
//	m, e := cmd.Result()
//	if e != nil {
//		return nil, e
//	}
//	return m, nil
//}
//
//func DelRoomUser(ctx context.Context, roomId string, clientName string) {
//	redis.Rdb.SRem(ctx, RoomMembers+roomId, clientName)
//}

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
