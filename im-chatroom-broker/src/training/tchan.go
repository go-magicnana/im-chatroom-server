package training

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func Gochan() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	for i := 1000; i < 1500; i++ {

		go func(i int) {
			ch := make(chan string, 1)
			go func() {
				for {
					select {
					case data := <-ch:
						fmt.Println(data)
					default:
						continue
					}
				}
			}()

			go func(i int) {
				index := 0
				for {
					index++
					ch <- strconv.Itoa(i)+" haha" + strconv.Itoa(index)
					time.Sleep(time.Second)
				}
			}(i)
		}(i)

	}

	wg.Wait()
	fmt.Println("OK")
}
