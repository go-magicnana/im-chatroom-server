package server

import (
	"im-chatroom-broker/context"
	"time"
)

//HeartBeating, determine if client send a message within set time by GravelChannel
// 心跳计时，根据GravelChannel判断Client是否在设定时间内发来信息

func HeartBeating(context *context.Context, readerChannel chan byte, timeout int) {
	select {
	case _ = <-readerChannel:
		//Log(context.Conn.RemoteAddr().String(), "get message, keeping heartbeating...")
		context.Conn().SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		break
	case <-time.After(time.Second * 5):
		//Log("It's really weird to get Nothing!!!")
		context.Conn().Close()
	}

}

func GravelChannel(n []byte, mess chan byte) {
	for _, v := range n {
		mess <- v
	}
	close(nil,nil,nil)
}
