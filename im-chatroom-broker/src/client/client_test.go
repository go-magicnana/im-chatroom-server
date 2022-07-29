package client

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

func TestStart(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go Start("send", "127.0.0.1")

	for i := 0; i < 2; i++ {
		go Start("Receiver", "127.0.0.1")
	}
	wg.Wait()

}

func TestStart2(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go Start("send", "47.95.148.3")

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
		go Start("Receiver", "47.95.148.3")
	}
	wg.Wait()

}

func Test_writeFile(t *testing.T) {
	path := "/Users/jinsong/work/TEST.txt"
	info := "hello"

	fi, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fi.Close()

	_,err = fi.WriteString(info)
	if err != nil {
		return
	}
}
