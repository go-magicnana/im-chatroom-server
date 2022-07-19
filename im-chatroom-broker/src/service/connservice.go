package service

import (
	"im-chatroom-broker/context"
	"im-chatroom-broker/protocol"
	"sync"
)

var users sync.Map
var dirty sync.Map

func SetUserContext(user *protocol.UserDevice, c *context.Context) {
	users.Store(user.UserKey, c)
}

func DelUserContext(userKey string) {
	users.Delete(userKey)
}

func GetUserContext(userKey string) (*context.Context, bool) {
	v, e := users.Load(userKey)

	return v.(*context.Context), e
}

func RangeUserContextAll(f func(key, value any) bool){
	users.Range(f)
}
