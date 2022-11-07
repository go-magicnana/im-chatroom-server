package client

import (
	"fmt"
	"golang.org/x/net/context"
	"im-chatroom-broker/service"
	"sync"
	"testing"
	"time"
)

var server = "192.168.3.242"

//var server = "127.0.0.1"

var user = &userInClient{
	userId: "1001",
	token:  "dltq",
	roomId: "100",
}

func TestRead(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	for i := 0; i < 2000; i++ {
		go Start("Receiver", server, 10, i, user)

		//if banch>=20 {
		//	time.Sleep(time.Second*3)
		//	banch = 0
		//}
	}

	wg.Wait()

}

func TestWrite(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	//go Start("send", "127.0.0.1","dltq","100",50)
	go Start("send", server, 200, 0, user)

	wg.Wait()

}

func TestCard(t *testing.T) {
	fmt.Println(service.CardMessage(context.Background()))
}

func TestReadMulti(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	//go Start("send", "47.95.148.3")

	for j := 0; j < 2000; j++ {
		//go Start("Receiver", "127.0.0.1","dltq","100",10)
		go Start("Receiver", server, 10, j, user)
		time.Sleep(time.Microsecond * 100)
	}

	//go Start("send", "127.0.0.1","dltq","100",50)
	//go Start("send", "47.95.148.3")

	wg.Wait()

}
