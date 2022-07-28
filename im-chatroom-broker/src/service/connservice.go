package service

import (
	"im-chatroom-broker/context"
	"sync"
)

var users sync.Map
var rooms sync.Map

func SetUserContext(clientName string, c *context.Context) {
	users.Store(clientName, c)
}

func DelUserContext(clientName string) (*context.Context, bool) {
	c, f := GetUserContext(clientName)
	users.Delete(clientName)
	return c, f
}

func GetUserContext(clientName string) (*context.Context, bool) {
	v, e := users.Load(clientName)

	if e {
		return v.(*context.Context), e
	} else {
		return nil, e
	}
}

func RangeUserContextAll(f func(key, value any) bool) {
	users.Range(f)
}

func SetRoomClient(roomId,clientName,userId string) {
	v, b := rooms.Load(roomId)
	if !b || v == nil {
		var m *sync.Map = new(sync.Map)
		rooms.Store(roomId,m)
		v,_ = rooms.Load(roomId)
	}
	v.(*sync.Map).Store(clientName,userId)
}

func RangeRoom(roomId string,f func(key, value any) bool) {
	v, b := rooms.Load(roomId)

	if !b || v == nil {
		return
	} else {
		v.(*sync.Map).Range(f)
	}
}

func GetRoomMap(roomId string) *sync.Map{
	a,b:=rooms.Load(roomId)

	if !b || a==nil {
		return nil
	}else{
		return a.(*sync.Map)
	}

}

func DelRoomClients(roomId,clientName string){
	m := GetRoomMap(roomId)
	if m!=nil{
		m.Delete(clientName)
	}
}

