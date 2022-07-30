package thread

import (
	"fmt"
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"go.uber.org/atomic"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/zaplog"
	"net"
	"sync"
)

var t sync.Map
var r sync.Map
var roomLock sync.Mutex

var Count *atomic.Int32

func init() {
	Count = atomic.NewInt32(0)
}

type ConnectClient struct {
	Broker     string
	ClientName string
	UserId     string
	RoomId     string
	Channel    chan *protocol.InnerPacket
	Conn       net.Conn
}

func (cc *ConnectClient) ToString() string {
	return cc.ClientName + " " + cc.UserId + " " + cc.RoomId + " " + fmt.Sprintf("%p", cc.Channel) + " " + cc.Conn.RemoteAddr().String()
}

func SetChannel(clientName string, cc *ConnectClient) {
	t.Store(clientName, cc)
	Count.Inc()
	zaplog.Logger.Debugf("ThreadContext SetChannel %s", cc.ToString())

}

func SetRoomChannels(roomId string, clientName string) {
	roomLock.Lock()
	v, _ := r.Load(roomId)

	if v == nil {
		m := linkedhashmap.New()
		m.Put(clientName, true)
		v = m
		r.Store(roomId, m)
	} else {
		l := v.(*linkedhashmap.Map)
		l.Put(clientName, true)
	}
	roomLock.Unlock()

}

func GetChannel(clientName string) *ConnectClient {
	k, b := t.Load(clientName)
	zaplog.Logger.Debugf("ThreadContext GetChannel %s %v %v", clientName, k, b)

	if k == nil {
		return nil
	}

	r := k.(*ConnectClient)

	return r

}

func GetRoomChannels(roomId string) []interface{} {
	v, _ := r.Load(roomId)

	if v == nil {
		return nil
	} else {
		temp := v.(*linkedhashmap.Map)
		return temp.Keys()
	}

}

func RemRoomChannel(roomId, clientName string) {
	v, _ := r.Load(roomId)
	if v != nil {
		temp := v.(*linkedhashmap.Map)
		temp.Remove(clientName)
	}
}

func RemChannel(clientName string) {
	t.Delete(clientName)
	zaplog.Logger.Debugf("ThreadContext RemChannel %s", clientName)

	Count.Dec()

}

func RanChannel(f func(key, value any) bool) {
	t.Range(f)
}
