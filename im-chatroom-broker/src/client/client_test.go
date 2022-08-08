package client

import (
	"sync"
	"testing"
	"time"
)

var server = "192.168.3.242"
//var server = "127.0.0.1"

func TestRead(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)


	banch := 0

	for i := 0; i < 2000; i++ {
		go Start("Receiver", server, "dltq", "100", 10, i)
		banch++

		if banch>=20 {
			time.Sleep(time.Second)
			banch = 0
		}
	}

	wg.Wait()

}

func TestWrite(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	//go Start("send", "127.0.0.1","dltq","100",50)
	go Start("send", server, "dltq", "100", 50, 0)

	wg.Wait()

}

func TestReadMulti(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	//go Start("send", "47.95.148.3")

	for j := 0; j < 2000; j++ {
		//go Start("Receiver", "127.0.0.1","dltq","100",10)
		go Start("Receiver", server, "dltq", "100", 10, j)
		time.Sleep(time.Microsecond * 100)
	}

	//go Start("send", "127.0.0.1","dltq","100",50)
	//go Start("send", "47.95.148.3")

	wg.Wait()

}
