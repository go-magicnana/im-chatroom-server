package service

import (
	"golang.org/x/net/context"
	"im-chatroom-gateway/redis"
)

const (
	RoomInfo string = "imchatroom:room.info:"
	RoomBlacks   string = "imchatroom:room.blacks:"
	RoomInstance string = "imchatroom:room.instance"
	RoomClients string = "imchatroom:room.clients:"

)

func SetRoomInstance(ctx context.Context, roomId string) int64 {
	return redis.Rdb.SAdd(ctx, RoomInstance, roomId).Val()
}

func GetRoomInstances(ctx context.Context) []string{
	return redis.Rdb.SMembers(ctx,RoomInstance).Val()
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

func RemRoomClients(broker, roomId, ClientName string) int64 {
	return redis.Rdb.SRem(context.Background(), RoomClients+broker+":"+roomId, ClientName).Val()
}

func DelRoomClients(broker, roomId string) int64 {
	return redis.Rdb.Del(context.Background(), RoomClients+broker+":"+roomId).Val()
}
