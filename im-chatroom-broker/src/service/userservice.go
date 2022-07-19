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

	/*hash*/
	UserDevice string = "imchatroom:userdevice:"
	UserInfo   string = "imchatroom:userinfo:"

	/*string json */
	UserAuth string = "imchatroom:userauth:"

	UserClients string = "imchatroom:userclients:"
)

func SetUserClient(ctx context.Context, userId string, userKey string) int64 {
	rdb := redis.Singleton()

	cmd := rdb.HSet(ctx, UserClients+userId, userKey, util.CurrentSecond())

	if cmd == nil {
		return 0
	}

	return cmd.Val()
}

func GetUserClients(ctx context.Context, userId string) []string {
	rdb := redis.Singleton()
	cmd := rdb.HGetAll(ctx, UserClients+userId)

	m := cmd.Val()

	ret := make([]string, 0)

	for k, _ := range m {
		ret = append(ret, k)
	}

	return ret
}

func GetUserAuth(ctx context.Context, token string) (*protocol.UserAuth, error) {
	redis := redis.Singleton()
	cmd := redis.Get(ctx, UserAuth+token)

	bs, err := cmd.Bytes()
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

func SetUserInfo(ctx context.Context, info protocol.UserInfo) {
	redis := redis.Singleton()

	bs, e := json.Marshal(info)

	if e != nil {
		util.Panic(e)
	}

	json := string(bs)

	redis.Set(ctx, UserInfo+info.UserId, json, time.Minute)
}

func GetUserInfo(ctx context.Context, userId string) (*protocol.UserInfo, error) {
	redis := redis.Singleton()

	cmd := redis.Get(ctx, UserInfo+userId)

	bs, err := cmd.Bytes()
	if err != nil {
		return nil, err
	}

	user := &protocol.UserInfo{}
	e2 := json.Unmarshal(bs, user)
	return user, e2
}

//func SetUserAlive(ctx context.Context, userKey string) {
//	redis := redis.Singleton()
//	redis.Expire(ctx, UserDevice+userKey, time.Second*20)
//}
//func DelUserRoom(ctx context.Context, userKey string) {
//	redis := redis.Singleton()
//	redis.HDel(ctx, UserDevice+userKey, "roomId")
//}

/**
UserKey string `json:"userKey"`
UserId  string `json:"userId"`
Device  string `json:"device"`
State   string `json:"state"`
RoomId  string `json:"roomId"`
Broker  string `json:"broker"`
*/

func SetUserDevice(ctx context.Context, user protocol.UserDevice) {

	redis := redis.Singleton()

	if util.IsNotEmpty(user.UserId) {
		redis.HSet(ctx, UserDevice+user.UserKey, "userId", user.UserId)
	}

	if util.IsNotEmpty(user.Broker) {
		redis.HSet(ctx, UserDevice+user.UserKey, "broker", user.Broker)
	}

	if util.IsNotEmpty(user.Device) {
		redis.HSet(ctx, UserDevice+user.UserKey, "device", user.Device)
	}

	if util.IsNotEmpty(user.RoomId) {
		redis.HSet(ctx, UserDevice+user.UserKey, "roomId", user.RoomId)
	}

	if util.IsNotEmpty(user.State) {
		redis.HSet(ctx, UserDevice+user.UserKey, "state", user.State)
	}

}

func GetUserDevice(ctx context.Context, userKey string) (*protocol.UserDevice, error) {

	redis := redis.Singleton()

	cmd := redis.HGetAll(ctx, UserDevice+userKey)
	m := cmd.Val()

	userDevice := &protocol.UserDevice{
		UserKey: userKey,
		UserId:  m["userId"],
		Device:  m["device"],
		State:   m["state"],
		RoomId:  m["roomId"],
		Broker:  m["broker"],
	}

	return userDevice, nil

}

func GetUserDeviceBroker(ctx context.Context, userKey string) (string, error) {
	redis := redis.Singleton()
	cmd := redis.HGet(ctx, UserDevice+userKey, "broker")
	return cmd.Val(), nil
}

func SetUserDevice2InRoom(ctx context.Context, userKey, roomId string) {
	redis := redis.Singleton()
	if util.IsNotEmpty(roomId) {
		redis.HSet(ctx, UserDevice+userKey, "roomId", roomId)
	}
}

func DelUserDeviceInRoom(ctx context.Context, userKey string) {
	redis := redis.Singleton()
	redis.HDel(ctx, UserDevice+userKey, "roomId")

}

func SetUserDevice2Login(ctx context.Context, userKey string, state int32) {
	redis := redis.Singleton()
	redis.HSet(ctx, UserDevice+userKey, "state", state)

}

func DelUserInfo(ctx context.Context, userKey string) {
	redis := redis.Singleton()
	redis.Del(ctx, UserInfo+userKey)
}
