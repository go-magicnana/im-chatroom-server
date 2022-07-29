package client

import (
	"sync"
	"testing"
)

func TestStart(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go Start("send","47.95.148.3")

	//for i := 0; i < 200; i++ {
	//	go Start("Receiver","47.95.148.3")
	//}
	wg.Wait()

}

func TestStart2(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go Start("send","47.95.148.3")

	//for i := 0; i < 200; i++ {
	//	go Start("Receiver","47.95.148.3")
	//}
	wg.Wait()

}

func TestStart3(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	//go Start("send","47.95.148.3")

	for i := 0; i < 1; i++ {
		go Start("Receiver","47.95.148.3")
	}
	wg.Wait()

}


