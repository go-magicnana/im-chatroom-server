package service

import (
	"im-chatroom-broker/context"
	"im-chatroom-broker/protocol"
	"sync"
)

var users sync.Map
var dirty sync.Map

func SetUserContext(user *protocol.UserDevice, c *context.Context) {
	users.Store(user.ClientName, c)
}

func DelUserContext(clientName string) (*context.Context, bool) {
	c, f := GetUserContext(clientName)
	users.Delete(clientName)
	return c, f
}

func GetUserContext(clientName string) (*context.Context, bool) {
	v, e := users.Load(clientName)

	if e {
		return v.(*context.Context), e
	} else {
		return nil, e
	}
}

func RangeUserContextAll(f func(key, value any) bool) {
	users.Range(f)
}
