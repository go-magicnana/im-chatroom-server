package context

import (
	"context"
	"im-chatroom-gateway/src/util"
	"net"
)

type Context struct {
	RemoteAddr  string
	BrokerAddr  string
	Conn        net.Conn
	Ctx         context.Context
	CancelFunc  context.CancelFunc
	ConnectTime int64
}

func (c Context) Close() {
	c.Conn.Close()
}

func NewContext(remoteAddr string, brokerAddr string, conn net.Conn, ctx context.Context, cancelFunc context.CancelFunc) *Context {
	return &Context{
		RemoteAddr:  remoteAddr,
		BrokerAddr:  brokerAddr,
		Conn:        conn,
		Ctx:         ctx,
		CancelFunc:  cancelFunc,
		ConnectTime: util.CurrentMillionSecond(),
	}
}
