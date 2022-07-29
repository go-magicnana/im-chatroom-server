package context

import (
	"go.uber.org/atomic"
	"im-chatroom-broker/protocol"
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
	clientName string
	userId     string
	roomId     string
	broker     string
	conn       net.Conn
	state      *atomic.Int32
	pingTime   *atomic.Int64
	readable   *atomic.Bool
	writeable  *atomic.Bool
	channel    chan *protocol.InnerPacket
}

func NewContext(brokerAddr string, conn net.Conn) *Context {

	return &Context{
		broker:    brokerAddr,
		conn:      conn,
		state:     atomic.NewInt32(Ready),
		pingTime:  atomic.NewInt64(0),
		readable:  atomic.NewBool(false),
		writeable: atomic.NewBool(false),
		channel:   make(chan *protocol.InnerPacket, 65535),
	}
}

func (c *Context) Readable() {
	c.readable.Store(true)
}

func (c *Context) UnReadable() {
	c.readable.Store(false)
}

func (c *Context) Writeable() {
	c.writeable.Store(true)
}

func (c *Context) Write(p *protocol.InnerPacket) {
	c.channel <- p
}

func (c *Context) Read() <- chan *protocol.InnerPacket {
	return c.channel
}

func (c *Context) ClientName() string {
	return c.clientName
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

func (c *Context) Connect(clientName string) (int32, bool) {
	ret := c.state.CAS(Ready, Connected)
	if ret {
		c.clientName = clientName
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

func (c *Context) Login(userId string) (int32, bool) {
	ret := c.state.CAS(Connected, Login)
	if ret {
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
