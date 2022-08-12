package ctx

import (
	"github.com/panjf2000/gnet/v2"
	"go.uber.org/atomic"
	"im-chatroom-broker/zaplog"
	"sync"
)

var conn sync.Map

var connCount *atomic.Int32

var BrokerAddress string

func init() {
	connCount = atomic.NewInt32(0)
}

type Context struct {
	Broker       string
	ClientName   string
	UserId       string
	RoomId       string
	Conn         gnet.Conn
	Time         int64
}

func (cc *Context) ToString() string {
	return cc.ClientName + " " + cc.UserId + " " + cc.RoomId + " " + cc.Conn.RemoteAddr().String()
}

func OpenContext(clientName string, cc *Context) {
	conn.Store(clientName, cc)
	connCount.Inc()
}

func GetContext(clientName string) *Context {
	k, b := conn.Load(clientName)

	if k == nil {
		zaplog.Logger.Debugf("ThreadContext GetChannel NotExist %s %v", clientName, b)
		return nil
	}

	r := k.(*Context)

	return r

}

func RemContext(clientName string) {
	_, f := conn.LoadAndDelete(clientName)

	if f {
		connCount.Dec()
	}

}

func RangeContext(f func(key, value any) bool) {
	conn.Range(f)
}

func ConnCount() int32{
	return connCount.Load()
}
