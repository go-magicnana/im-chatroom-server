package handler

import (
	"im-chatroom-broker/context"
	"im-chatroom-broker/protocol"
	"sync"
)

var users sync.Map
var dirty sync.Map



func SetUserContext(user *protocol.User, c *context.Context) {
	users.Store(user.UserId, c)
}

func GetUserContext(userId string) *context.Context {
	v, f := users.Load(userId)

	if f {
		return v.(*context.Context)
	} else {
		return nil
	}

}

//TODO ... 如果超过一段时间没有connect，就关了他
