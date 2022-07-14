package training

import (
	"fmt"
	"sync"
	"time"
)

func Gochan() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	ch := make(chan byte, 1)
	ch2 := make(chan byte,1)

	go func() {
		ret := <-ch //如果没有内容，则阻塞
		fmt.Println("read from chan", ret)
		wg.Done()

		rr := <- ch2
		fmt.Println(rr)
	}()

	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second)
		}
		ch <- 1
	}()

	wg.Wait()
	fmt.Println("OK")
}
