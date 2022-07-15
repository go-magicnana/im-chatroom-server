package rediss

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSingleton(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		RedisSingleton().Set(context.Background(), "test1", "hahs", time.Hour)
		wg.Done()
	}()
	go func() {
		RedisSingleton().HSet(context.Background(), "test2", "field", "haha")
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("OK")
}
