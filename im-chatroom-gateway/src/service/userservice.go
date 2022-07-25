package service

import (
	"encoding/json"
	"golang.org/x/net/context"
	"im-chatroom-gateway/protocol"
	"im-chatroom-gateway/redis"
)

const (

	/*hash*/
	UserDevice string = "imchatroom:user.device:"
	UserInfo   string = "imchatroom:user.info:"

	/*string json */
	UserAuth string = "imchatroom:user.auth:"

	UserClients string = "imchatroom:user.clients:"
)

func GetUserClients(ctx context.Context, userId string) []string {
	cmd := redis.Rdb.HGetAll(ctx, UserClients+userId)

	m := cmd.Val()

	ret := make([]string, 0)

	for k, _ := range m {
		ret = append(ret, k)
	}

	return ret
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

func GetUserDeviceBroker(ctx context.Context, clientName string) (string, error) {
	cmd := redis.Rdb.HGet(ctx, UserDevice+clientName, "broker")
	return cmd.Val(), nil
}
