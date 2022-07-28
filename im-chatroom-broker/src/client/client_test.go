package client

import (
	"sync"
	"testing"
)

func TestStart(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go Start("send","127.0.0.1")

	for i := 0; i < 100; i++ {
		go Start("Receiver","127.0.0.1")
	}
	wg.Wait()

}

