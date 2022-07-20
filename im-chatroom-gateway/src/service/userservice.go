package service

import (
	"encoding/json"
	"golang.org/x/net/context"
	"im-chatroom-gateway/src/protocol"
	"im-chatroom-gateway/src/redis"
	"im-chatroom-gateway/src/util"
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
	cmd := redis.Rdb.HSet(ctx, UserClients+userId, userKey, util.CurrentSecond())

	if cmd == nil {
		return 0
	}

	return cmd.Val()
}

func GetUserClients(ctx context.Context, userId string) []string {
	cmd := redis.Rdb.HGetAll(ctx, UserClients+userId)

	m := cmd.Val()

	ret := make([]string, 0)

	for k, _ := range m {
		ret = append(ret, k)
	}

	return ret
}

func GetUserAuth(ctx context.Context, token string) (*protocol.UserAuth, error) {
	cmd := redis.Rdb.Get(ctx, UserAuth+token)

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

func DelUserAuth(ctx context.Context, token string) {
	//redis := redis.Rdb
	//redis.Del(ctx, UserAuth+token)
}

func SetUserInfo(ctx context.Context, info protocol.UserInfo) {
	bs, e := json.Marshal(info)

	if e != nil {
		util.Panic(e)
	}

	json := string(bs)

	redis.Rdb.Set(ctx, UserInfo+info.UserId, json, time.Minute)
}

func GetUserInfo(ctx context.Context, userId string) (*protocol.UserInfo, error) {
	cmd := redis.Rdb.Get(ctx, UserInfo+userId)

	bs, err := cmd.Bytes()
	if err != nil {
		return nil, err
	}

	if len(bs) == 0 {
		return nil, nil
	}

	user := &protocol.UserInfo{}
	e2 := json.Unmarshal(bs, user)
	return user, e2
}

func SetUserAlive(ctx context.Context, userId, userKey string) {
	redis.Rdb.Expire(ctx, UserDevice+userKey, time.Second*20)
	redis.Rdb.Expire(ctx, UserInfo+userId, time.Second*20)
	redis.Rdb.Expire(ctx, UserClients+userId, time.Second*20)
}

func SetUserDevice(ctx context.Context, user protocol.UserDevice) {

	if util.IsNotEmpty(user.UserId) {
		redis.Rdb.HSet(ctx, UserDevice+user.UserKey, "userId", user.UserId)
	}

	if util.IsNotEmpty(user.Broker) {
		redis.Rdb.HSet(ctx, UserDevice+user.UserKey, "broker", user.Broker)
	}

	if util.IsNotEmpty(user.Device) {
		redis.Rdb.HSet(ctx, UserDevice+user.UserKey, "device", user.Device)
	}

	if util.IsNotEmpty(user.RoomId) {
		redis.Rdb.HSet(ctx, UserDevice+user.UserKey, "roomId", user.RoomId)
	}

	if util.IsNotEmpty(user.State) {
		redis.Rdb.HSet(ctx, UserDevice+user.UserKey, "state", user.State)
	}

}

func GetUserDevice(ctx context.Context, userKey string) (*protocol.UserDevice, error) {

	cmd := redis.Rdb.HGetAll(ctx, UserDevice+userKey)
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
	cmd := redis.Rdb.HGet(ctx, UserDevice+userKey, "broker")
	return cmd.Val(), nil
}

func SetUserDevice2InRoom(ctx context.Context, userKey, roomId string) {
	if util.IsNotEmpty(roomId) {
		redis.Rdb.HSet(ctx, UserDevice+userKey, "roomId", roomId)
	}
}

func DelUserDeviceInRoom(ctx context.Context, userKey string) {
	redis.Rdb.HDel(ctx, UserDevice+userKey, "roomId")
}

func DelUserDevice(ctx context.Context, userKey string) {
	redis.Rdb.Del(ctx, UserDevice+userKey)
}

func SetUserDevice2Login(ctx context.Context, userKey string, state int32) {
	redis.Rdb.HSet(ctx, UserDevice+userKey, "state", state)

}

func DelUserInfo(ctx context.Context, userKey string) {
	redis.Rdb.Del(ctx, UserInfo+userKey)
}
