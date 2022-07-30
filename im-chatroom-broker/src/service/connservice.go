package service

//var users sync.Map
//var rooms sync.Map
//
//var _chOnce sync.Once
//
//var CH *ContextHolder
//
//func SingletonContextHolder() *ContextHolder {
//	_chOnce.Do(func() {
//		CH = newContextHolder()
//	})
//
//	return CH
//}
//
//type ContextHolder struct {
//	channelMap sync.Map
//}
//
//func newContextHolder() *ContextHolder {
//	return &ContextHolder{}
//}

//func SetUserContext(clientName string, c *context.Context) {
//	users.Store(clientName, c)
//}
//
//func DelUserContext(clientName string) (*context.Context, bool) {
//	c, f := GetUserContext(clientName)
//	users.Delete(clientName)
//	return c, f
//}
//
//func GetUserContext(clientName string) (*context.Context, bool) {
//	v, e := users.Load(clientName)
//
//	if e {
//		return v.(*context.Context), e
//	} else {
//		return nil, e
//	}
//}
//
//func RangeUserContextAll(f func(key, value any) bool) {
//	users.Range(f)
//}
//
//func SetRoomClient(roomId, clientName, userId string) {
//	v, b := rooms.Load(roomId)
//	if !b || v == nil {
//		var m *sync.Map = new(sync.Map)
//		rooms.Store(roomId, m)
//		v, _ = rooms.Load(roomId)
//	}
//	v.(*sync.Map).Store(clientName, userId)
//}
//
//func RangeRoom(roomId string, f func(key, value any) bool) {
//	v, b := rooms.Load(roomId)
//
//	if !b || v == nil {
//		return
//	} else {
//		v.(*sync.Map).Range(f)
//	}
//}
//
//func GetRoomMap(roomId string) *sync.Map {
//	a, b := rooms.Load(roomId)
//
//	if !b || a == nil {
//		return nil
//	} else {
//		return a.(*sync.Map)
//	}
//
//}

//func DelRoomClients(roomId, clientName string) {
//	m := GetRoomMap(roomId)
//	if m != nil {
//		m.Delete(clientName)
//	}
//}

//func (ch *ContextHolder) SetChannel(clientName string, channel chan *protocol.InnerPacket) {
//	ch.channelMap.Store(clientName, channel)
//}
//
//func (ch *ContextHolder) GetChannel(clientName string) chan *protocol.InnerPacket {
//	k, _ := ch.channelMap.Load(clientName)
//
//	return k.(chan *protocol.InnerPacket)
//
//}
//
//func (ch *ContextHolder) RemChannel(clientName string) {
//	ch.channelMap.Delete(clientName)
//}
//
//func (ch *ContextHolder) RanChannel(f func(key, value any) bool) {
//	ch.channelMap.Range(f)
//}
