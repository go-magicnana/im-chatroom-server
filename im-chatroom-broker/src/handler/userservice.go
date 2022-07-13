package handler

import (
	"encoding/json"
	"golang.org/x/net/context"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/redis"
	"im-chatroom-broker/util"
	"time"
)

const (
	UserDetail string = "imchatroom:userdetail:"
	Broker     string = "imchatroom:broker:"
)

func SetUserInfo(ctx context.Context, user protocol.User) {

	redis := redis.Singleton()
	redis.Set(ctx, UserDetail+user.UserId, util.ToJsonString(user), time.Minute*1)

	if util.IsNotEmpty(user.Name) {
		redis.HSet(ctx, UserDetail+user.UserId, "name", user.Name)
	}

	if util.IsNotEmpty(user.Avatar) {
		redis.HSet(ctx, UserDetail+user.UserId, "avatar", user.Avatar)
	}

	if util.IsNotEmpty(user.Role) {
		redis.HSet(ctx, UserDetail+user.UserId, "role", user.Role)
	}

	if util.IsNotEmpty(user.Broker) {
		redis.HSet(ctx, UserDetail+user.UserId, "broker", user.Broker)
	}

	if util.IsNotEmpty(user.RoomId) {
		redis.HSet(ctx, UserDetail+user.UserId, "roomId", user.RoomId)
	}
}

func GetUserInfo(ctx context.Context,userId string) protocol.User {
	user := protocol.User{}
	cmd := redis.Singleton().Get(ctx, UserDetail+userId)
	bs, _ := cmd.Bytes()
	json.Unmarshal(bs, &user)
	return user
}
