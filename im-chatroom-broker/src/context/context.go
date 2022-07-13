package context

import (
	"im-chatroom-broker/util"
	"net"
)

type Context struct {
	UserId      string
	BrokerAddr  string
	Conn        net.Conn
	ConnectTime int64
}

func NewContext(brokerAddr string, conn net.Conn) *Context {
	return &Context{
		BrokerAddr:  brokerAddr,
		Conn:        conn,
		ConnectTime: util.CurrentMillionSecond(),
	}
}

func (c Context) PutUser(userId string)  {
	c.UserId = userId
}
