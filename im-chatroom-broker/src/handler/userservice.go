package handler

import (
	"encoding/json"
	"golang.org/x/net/context"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/redis"
	"im-chatroom-broker/util"
)

const (

	/*hash*/
	UserInfo string = "imchatroom:userinfo:"

	/*string json */
	UserAuth string = "imchatroom:userauth:"
	Broker   string = "imchatroom:broker:"
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

func SetUserInfo(ctx context.Context, user *protocol.User) {

	redis := redis.Singleton()

	if util.IsNotEmpty(user.Name) {
		redis.HSet(ctx, UserInfo+user.UserId, "name", user.Name)
	}

	if util.IsNotEmpty(user.Token){
		redis.HSet(ctx,UserInfo+user.UserId,"token",user.Token)
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

	if util.IsNotEmpty(user.RoomId) {
		redis.HSet(ctx, UserInfo+user.UserId, "roomId", user.RoomId)
	}
}

