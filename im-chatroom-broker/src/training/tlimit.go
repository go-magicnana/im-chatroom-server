package training

import (
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/time/rate"
	"time"
)

func limitAllow() {
	limiter := rate.NewLimiter(0, 2)

	for i := 100; i < 200; i++ {

		if limiter.Allow() {
			go func() {
				fmt.Println(i, " Do something ...")
				limiter.Wait(context.Background())
			}()
		} else {
			fmt.Println("No token")
		}

	}

	time.Sleep(time.Hour * 10)

}

func limitWait() {
	limiter := rate.NewLimiter(1, 2)

	for i := 100; i < 110; i++ {

		if err := limiter.Wait(context.Background()); err==nil {
			go func() {
				fmt.Println(i, " Do something ...")
				limiter.Wait(context.Background())
			}()
		}else{
			fmt.Println(err)
		}
		fmt.Println("continue")
	}

	time.Sleep(time.Hour * 10)
}
