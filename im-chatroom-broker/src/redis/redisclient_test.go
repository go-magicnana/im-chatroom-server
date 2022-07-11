package redis

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
		Singleton().Set(context.Background(),"test1","hahs",time.Hour)
		wg.Done()
	}()
	go func() {
		Singleton().HSet(context.Background(),"test2","field","haha")
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("OK")
}
