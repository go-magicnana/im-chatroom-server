package training

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}


func root() {

	fmt.Println("root start")

	wg.Add(1)

	root := context.Background()
	c, cc := context.WithCancel(root)

	go child(c, cc)

	go mama(c, cc)

	wg.Wait()
}

func child(ctx context.Context, cancel context.CancelFunc) {
	i := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Println("child done")
			wg.Done()
			return
		default:
			i++
			time.Sleep(time.Second * 2)
			fmt.Println("child do ", i)
		}
	}
	wg.Done()
}

func mama(ctx context.Context, cancel context.CancelFunc) {
	time.Sleep(time.Second*10)
	cancel()
}
