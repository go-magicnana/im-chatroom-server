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

	/*hash*/
	UserInfo string = "imchatroom:userinfo:"

	/*string json */
	UserAuth string = "imchatroom:userauth:"
)

func GetUserAuth(ctx context.Context, token string) (*protocol.User, error) {
	redis := redis.Singleton()
	cmd := redis.Get(ctx, UserAuth+token)

	bs, err := cmd.Bytes()
	if err != nil {
		return nil, err
	} else {
		user := &protocol.User{}
		e := json.Unmarshal(bs, user)
		if e != nil {
			return nil, e
		} else {
			return user, nil
		}
	}
}

func SetUserRoom(ctx context.Context, userKey, roomId string) {
	redis := redis.Singleton()
	if util.IsNotEmpty(roomId) {
		redis.HSet(ctx, UserInfo+userKey, "roomId", roomId)
		redis.Expire(ctx, UserInfo+userKey, time.Second*20)
	}
}

func SetUserLogin(ctx context.Context, userKey string, state int32) {
	redis := redis.Singleton()
	redis.HSet(ctx, UserInfo+userKey, "state", state)
	redis.Expire(ctx, UserInfo+userKey, time.Second*20)

}

func SetUserAlive(ctx context.Context, userKey string) {
	redis := redis.Singleton()
	redis.Expire(ctx, UserInfo+userKey, time.Second*20)
}
func DelUserRoom(ctx context.Context, userKey string) {
	redis := redis.Singleton()
	redis.HDel(ctx, UserInfo+userKey, "roomId")
}

func SetUserInfo(ctx context.Context, user *protocol.User) {

	redis := redis.Singleton()

	if util.IsNotEmpty(user.UserId) {
		redis.HSet(ctx, UserInfo+user.UserKey, "userId", user.UserId)
	}

	if util.IsNotEmpty(user.Name) {
		redis.HSet(ctx, UserInfo+user.UserKey, "name", user.Name)
	}

	if util.IsNotEmpty(user.Token) {
		redis.HSet(ctx, UserInfo+user.UserKey, "token", user.Token)
	}

	if util.IsNotEmpty(user.Avatar) {
		redis.HSet(ctx, UserInfo+user.UserKey, "avatar", user.Avatar)
	}

	if util.IsNotEmpty(user.Role) {
		redis.HSet(ctx, UserInfo+user.UserKey, "role", user.Role)
	}

	if util.IsNotEmpty(user.Broker) {
		redis.HSet(ctx, UserInfo+user.UserKey, "broker", user.Broker)
	}

	redis.Expire(ctx, UserInfo+user.UserKey, time.Second*20)

}

func DelUserInfo(ctx context.Context, userKey string) {
	redis := redis.Singleton()
	redis.Del(ctx, UserInfo+userKey)
}

func GetUserInfo(ctx context.Context, userKey string) (*protocol.User, error) {
	redis := redis.Singleton()
	cmd := redis.HGetAll(ctx, UserInfo+userKey)

	m, e := cmd.Result()

	if e != nil {
		return nil, e
	}
	user := &protocol.User{}
	user.RoomId = m["roomId"]
	user.UserId = m["userId"]
	user.UserKey = userKey
	user.Role = m["role"]
	user.Token = m["token"]
	user.Broker = m["broker"]
	user.Name = m["name"]
	user.Avatar = m["avatar"]
	user.State = m["state"]

	return user, nil

}
