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

func SetUserRoom(ctx context.Context, userId, roomId string) {
	redis := redis.Singleton()
	if util.IsNotEmpty(roomId) {
		redis.HSet(ctx, UserInfo+userId, "roomId", roomId)
	}
}

func SetUserLogin(ctx context.Context, userId string, state int32) {
	redis := redis.Singleton()
	redis.HSet(ctx, UserInfo+userId, "state", state)
}

func SetUserAlive(ctx context.Context,userId string){
	redis := redis.Singleton()
	redis.Expire(ctx, UserInfo+userId, time.Second*20)
}
func DelUserRoom(ctx context.Context, userId string) {
	redis := redis.Singleton()
	redis.HDel(ctx, UserInfo+userId, "roomId")
}

func SetUserInfo(ctx context.Context, user *protocol.User) {

	redis := redis.Singleton()

	if util.IsNotEmpty(user.Name) {
		redis.HSet(ctx, UserInfo+user.UserId, "name", user.Name)
	}

	if util.IsNotEmpty(user.Token) {
		redis.HSet(ctx, UserInfo+user.UserId, "token", user.Token)
	}

	if util.IsNotEmpty(user.Avatar) {
		redis.HSet(ctx, UserInfo+user.UserId, "avatar", user.Avatar)
	}

	if util.IsNotEmpty(user.Role) {
		redis.HSet(ctx, UserInfo+user.UserId, "role", user.Role)
	}

	if util.IsNotEmpty(user.Broker) {
		redis.HSet(ctx, UserInfo+user.UserId, "broker", user.Broker)
	}
}

func DelUserInfo(ctx context.Context, userId string) {
	redis := redis.Singleton()
	redis.Del(ctx, UserInfo+userId)
}

func GetUserInfo(ctx context.Context, userId string) (*protocol.User, error) {
	redis := redis.Singleton()
	cmd := redis.HGetAll(ctx, UserInfo+userId)

	m, e := cmd.Result()

	if e != nil {
		return nil, e
	}
	user := &protocol.User{}
	user.RoomId = m["roomId"]
	user.UserId = userId
	user.Role = m["role"]
	user.Token = m["token"]
	user.Broker = m["broker"]
	user.Name = m["name"]
	user.Avatar = m["avatar"]
	user.State = m["state"]

	return user, nil

}
