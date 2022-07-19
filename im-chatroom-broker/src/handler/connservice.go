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

func GetUserContext(userKey string) (*context.Context, bool) {
<<<<<<< HEAD
	v, e := users.Load(userKey)

	return v.(*context.Context), e
}

=======
	value, exist := users.Load(userKey)
	if exist {
		return value.(*context.Context), true
	} else {
		return nil, false
	}
}

//TODO ... 如果超过一段时间没有connect，就关了他
>>>>>>> ded2abccc3028f33a1e622ce850902b41be72c31
