package client

import (
	"sync"
	"testing"
)

func TestStart(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	for i := 0; i < 1; i++ {
		go Start("127.0.0.1")
	}
	wg.Wait()

}

func BenchmarkFib(b *testing.B) {
	var wg sync.WaitGroup
	wg.Add(1)
	for n := 0; n < b.N; n++ {
		go Start("127.0.0.1")
	}
	wg.Wait()
}