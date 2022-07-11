package server

import (
	"context"
	"im-chatroom-broker/util"
	"net"
)

type Context struct {
	RemoteAddr     string
	BrokerAddr     string
	Conn           net.Conn
	Ctx            context.Context
	CancelFunc     context.CancelFunc
	connectTime    int64
	authTime       int64
	lastActiveTime int64
}

func NewContext(remoteAddr string, brokerAddr string, conn net.Conn, ctx context.Context, cancelFunc context.CancelFunc) *Context {
	return &Context{
		RemoteAddr: remoteAddr,
		BrokerAddr: brokerAddr,
		Conn:       conn,
		Ctx:        ctx,
		CancelFunc: cancelFunc,
		connectTime: util.CurrentMillionSecond(),
	}
}
