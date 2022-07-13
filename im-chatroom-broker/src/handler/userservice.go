package handler

import (
	"encoding/json"
	"im-chatroom-broker/context"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/redis"
	"im-chatroom-broker/util"
	"sync"
	"time"
)

var users sync.Map
var dirty sync.Map

const (
	UserDetail string = "imchatroom:userdetail:"
	Broker     string = "imchatroom:broker:"
)

func SetUser(user protocol.User, c *context.Context) {

	redis := redis.Singleton()
	users.Store(user.UserId, c)
	redis.Set(c.Ctx, UserDetail+user.UserId, util.ToJsonString(user), time.Minute*1)

	if util.IsNotEmpty(user.Name) {
		redis.HSet(c.Ctx, UserDetail+user.UserId, "name", user.Name)
	}

	if util.IsNotEmpty(user.Avatar) {
		redis.HSet(c.Ctx, UserDetail+user.UserId, "avatar", user.Avatar)
	}

	if util.IsNotEmpty(user.Role) {
		redis.HSet(c.Ctx, UserDetail+user.UserId, "role", user.Role)
	}

	if util.IsNotEmpty(c.BrokerAddr) {
		redis.HSet(c.Ctx, UserDetail+user.UserId, "broker", c.BrokerAddr)
	}

	if util.IsNotEmpty(user.RoomId) {
		redis.HSet(c.Ctx, UserDetail+user.UserId, "roomId", user.RoomId)
	}
}

func GetUserConn(userId string) *context.Context {
	ay, exist := users.Load(userId)
	if exist {
		return ay.(*context.Context)
	} else {
		return nil
	}
}

func GetUserInfo(userId string, context *context.Context) protocol.User {
	user := protocol.User{}
	cmd := redis.Singleton().Get(context.Ctx, UserDetail+userId)
	bs, _ := cmd.Bytes()
	json.Unmarshal(bs, &user)
	return user
}
