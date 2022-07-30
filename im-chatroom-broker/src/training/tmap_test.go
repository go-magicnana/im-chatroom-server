package training

import (
	"fmt"
	"testing"
	"time"
)

func Test_go1(t *testing.T) {
	go1("go1",make(chan string,1))
	go1("go2",make(chan string,1))
	go1("go3",make(chan string,1))
	go1("go4",make(chan string,1))
	go1("go5",make(chan string,1))
	go1("go6",make(chan string,1))
	go1("go7",make(chan string,1))
	go1("go8",make(chan string,1))
	go1("go9",make(chan string,1))
	go1("go0",make(chan string,1))

	go func() {
		go2("go1")
	}()

	time.Sleep(time.Second)
	go3()
	time.Sleep(time.Second)
	fmt.Println("-----")
	go func() {
		go1("hi1",make(chan string,1))
		go1("hi2",make(chan string,1))
		go1("hi3",make(chan string,1))
	}()
	go3()
	time.Sleep(time.Second)	//
	fmt.Println("-----")

	go func() {
		go2("go1")
	}()
	//
	//time.Sleep(time.Second)
	//go3()
	//
	//
	//go func() {
	//	go2("hi1")
	//}()
	//
	//time.Sleep(time.Second*2)
	//
	go3()
	time.Sleep(time.Second)
}
