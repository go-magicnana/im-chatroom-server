package zaplog

import (
	"im-chatroom-broker/client"
	"sync"
	"testing"
)

func TestTest1(t *testing.T) {
	InitLogger()
	defer Logger.Sync()
	//Infof("这是一个描述 %s","binggo")
	//
	//Test1()

}

func TestFoo(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go client.Start("send","47.95.148.3")

	//for i := 0; i < 200; i++ {
	//	go Start("Receiver","47.95.148.3")
	//}
	wg.Wait()}
