package service

import (
	"encoding/json"
	"golang.org/x/net/context"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/redis"
	"im-chatroom-broker/util"
	"time"
)

const (
	UserInfo    string = "imchatroom:user.info:"
	UserAuth    string = "imchatroom:user.auth:"
	UserClients string = "imchatroom:user.clients:"
)

func SetUserClient(userId string, clientName string) int64 {
	return redis.Rdb.SAdd(context.Background(), UserClients+userId, clientName).Val()
}
func GetUserClients(userId string) []string {
	return redis.Rdb.SMembers(context.Background(), UserClients+userId).Val()

}
func RemUserClient(userId, clientName string) int64 {
	return redis.Rdb.SRem(context.Background(), UserClients+userId, clientName).Val()

}

func RefreshUserClient(userId string) {
	redis.Rdb.Expire(context.Background(), UserClients+userId, time.Minute*30)
}

func GetUserAuth(token string) (*protocol.UserAuth, error) {
	bs, err := redis.Rdb.Get(context.Background(), UserAuth+token).Bytes()
	if err != nil {
		return nil, err
	} else {
		user := &protocol.UserAuth{}
		e := json.Unmarshal(bs, user)
		if e != nil {
			return nil, e
		} else {
			return user, nil
		}
	}
}

func DelUserAuth(token string) {
	//redis := redis.Rdb
	//redis.Del(ctx, UserAuth+token)
}

func SetUserInfo(info protocol.UserInfo) string {
	bs, e := json.Marshal(info)

	if e != nil {
		util.Panic(e)
	}

	json := string(bs)

	return redis.Rdb.Set(context.Background(), UserInfo+info.UserId, json, time.Minute).Val()
}

func GetUserInfo(userId string) *protocol.UserInfo {
	cmd := redis.Rdb.Get(context.Background(), UserInfo+userId)

	bs, err := cmd.Bytes()
	if err != nil {
		return nil
	}

	if len(bs) == 0 {
		return nil
	}

	user := &protocol.UserInfo{}
	json.Unmarshal(bs, user)
	return user
}

func RefreshUserInfo(userId string) {
	redis.Rdb.Expire(context.Background(), UserInfo+userId, time.Minute*30)
}

//func SetUserAlive(ctx context.Context, userId, clientName string) {
//	redis.Rdb.Expire(ctx, UserDevice+clientName, time.Second*20)
//	redis.Rdb.Expire(ctx, UserInfo+userId, time.Second*20)
//	redis.Rdb.Expire(ctx, UserClients+userId, time.Second*20)
//}

//func SetUserDevice(ctx context.Context, user protocol.UserDevice) {
//
//	if util.IsNotEmpty(user.UserId) {
//		redis.Rdb.HSet(ctx, UserDevice+user.ClientName, "userId", user.UserId)
//	}
//
//	if util.IsNotEmpty(user.Broker) {
//		redis.Rdb.HSet(ctx, UserDevice+user.ClientName, "broker", user.Broker)
//	}
//
//	//if util.IsNotEmpty(user.Device) {
//	//	redis.Rdb.HSet(ctx, UserDevice+user.ClientName, "device", user.Device)
//	//}
//
//	if util.IsNotEmpty(user.RoomId) {
//		redis.Rdb.HSet(ctx, UserDevice+user.ClientName, "roomId", user.RoomId)
//	}
//
//	if util.IsNotEmpty(user.State) {
//		redis.Rdb.HSet(ctx, UserDevice+user.ClientName, "state", user.State)
//	}
//
//}
//
//func GetUserDevice(ctx context.Context, clientName string) (*protocol.UserDevice, error) {
//
//	cmd := redis.Rdb.HGetAll(ctx, UserDevice+clientName)
//	m := cmd.Val()
//
//	if m == nil || len(m) == 0 {
//		return nil, nil
//	}
//
//	userDevice := &protocol.UserDevice{
//		ClientName: clientName,
//		UserId:     m["userId"],
//		//Device:     m["device"],
//		State:      m["state"],
//		RoomId:     m["roomId"],
//		Broker:     m["broker"],
//	}
//
//	return userDevice, nil
//
//}
//
//func GetUserDeviceBroker(ctx context.Context, clientName string) (string, error) {
//	cmd := redis.Rdb.HGet(ctx, UserDevice+clientName, "broker")
//	return cmd.Val(), nil
//}
//
//func SetUserDevice2InRoom(ctx context.Context, clientName, roomId string) {
//	if util.IsNotEmpty(roomId) {
//		redis.Rdb.HSet(ctx, UserDevice+clientName, "roomId", roomId)
//	}
//}
//
//func DelUserDeviceInRoom(ctx context.Context, clientName string) {
//	redis.Rdb.HDel(ctx, UserDevice+clientName, "roomId")
//}
//
//func DelUserDevice(ctx context.Context, clientName string) {
//	redis.Rdb.Del(ctx, UserDevice+clientName)
//
//}
//
//func SetUserDevice2Login(ctx context.Context, clientName string, state int32) {
//	redis.Rdb.HSet(ctx, UserDevice+clientName, "state", state)
//}
