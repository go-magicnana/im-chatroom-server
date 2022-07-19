package context

import (
	"go.uber.org/atomic"
	"im-chatroom-broker/util"
	"net"
)

const (
	Closed    = -1
	Ready     = 0
	Connected = 1
	Login     = 2
	JoinRoom  = 3
)

type Context struct {
	userKey  string
	userId   string
	roomId   string
	broker   string
	conn     net.Conn
	state    *atomic.Int32
	pingTime *atomic.Int64
}

func NewContext(brokerAddr string, conn net.Conn) *Context {

	return &Context{
		broker:   brokerAddr,
		conn:     conn,
		state:    atomic.NewInt32(Ready),
		pingTime: atomic.NewInt64(0),
	}
}

func (c *Context) UserKey() string {
	return c.userKey
}

func (c *Context) UserId() string {
	return c.userId
}

func (c *Context) RoomId() string {
	return c.roomId
}

func (c *Context) Broker() string {
	return c.broker
}

func (c *Context) State() int32 {
	return c.state.Load()
}

func (c *Context) Conn() net.Conn {
	return c.conn
}

func (c *Context) Connect() (int32, bool) {
	ret := c.state.CAS(Ready, Connected)
	if ret {
		return Connected, true
	} else {
		return c.state.Load(), false
	}
}

func (c *Context) Ping() (int64, bool) {
	ret := util.CurrentSecond()
	c.pingTime.Store(ret)
	return ret, true
}

func (c *Context) Login(userKey, userId string) (int32, bool) {
	ret := c.state.CAS(Connected, Login)
	if ret {
		c.userKey = userKey
		c.userId = userId
		return Login, true
	} else {
		return c.state.Load(), false
	}
}

func (c *Context) JoinRoom(roomId string) (int32, bool) {
	ret := c.state.CAS(Login, JoinRoom)
	if ret {
		c.roomId = roomId
		return JoinRoom, true
	} else {
		return c.state.Load(), false
	}
}

func (c *Context) LeaveRoom() (int32, bool) {
	ret := c.state.CAS(JoinRoom, Login)
	if ret {
		c.roomId = ""
		return Login, true
	} else {
		return c.state.Load(), false
	}
}

func (c *Context) ChangeRoom(roomId string) {
	c.roomId = roomId
}

func (c *Context) Close() {
	c.state.Store(Closed)
	c.conn.Close()
}
