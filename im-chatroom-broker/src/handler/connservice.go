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

func DelUserContext(userId string){
	users.Delete(userId)
}


//TODO ... 如果超过一段时间没有connect，就关了他
