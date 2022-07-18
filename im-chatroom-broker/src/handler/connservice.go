package handler

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

//TODO ... 如果超过一段时间没有connect，就关了他
