package client

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	//for i := 0; i < 1; i++ {
	//	go Start("Receiver", "127.0.0.1")
	//	//go Start("Receiver", "47.95.148.3")
	//}

	//time.Sleep(time.Second*5)
	//go Start("send", "127.0.0.1")
	go Start("send", "47.95.148.3")

	wg.Wait()

}

func TestStart1(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	//go Start("send", "47.95.148.3")

	for i := 0; i < 1; i++ {
		for j := 0; j < 100; j++ {
			go Start("Receiver", "127.0.0.1")
		}
		time.Sleep(time.Second * 5)

		//go Start("Receiver", "127.0.0.1",strconv.Itoa(i))
	}
	wg.Wait()

}

func TestStart11(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	//go Start("send", "47.95.148.3")

	for i := 0; i < 10; i++ {
		for j := 0; j < 100; j++ {
			go Start("Receiver", "47.95.148.3")
		}
		time.Sleep(time.Second * 1)

		//go Start("Receiver", "127.0.0.1",strconv.Itoa(i))
	}
	wg.Wait()

}

func TestStart2(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	//go Start("send", "47.95.148.3")

	for i := 0; i < 300; i++ {
		go Start("Receiver", "47.95.148.3")
		//go Start("Receiver", "127.0.0.1",strconv.Itoa(i))
	}
	wg.Wait()

}

func TestStart3(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	//go Start("send", "47.95.148.3")

	for i := 0; i < 10; i++ {
		go Start("Receiver", "47.95.148.3")
		//go Start("Receiver", "127.0.0.1",strconv.Itoa(i))
		time.Sleep(time.Second * 5)
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

	_, err = fi.WriteString(info)
	if err != nil {
		return
	}
}

func TestSetUserAuth(t *testing.T) {
	for i := 1000; i < 2000; i++ {
		SetUserAuth(strconv.Itoa(i))
	}
}

func TestServer(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	for i := 0; i < 100; i++ {
		go thread(i)
	}
	wg.Wait()
}

func thread(i int) {
	c, _ := net.DialTimeout("tcp", "192.168.3.242:33121", 10*time.Second)
	time.Sleep(time.Millisecond)

	go func() {
		c.Write([]byte("haha"))
		time.Sleep(time.Second)
	}()

	go func() {

		bs := make([]byte,1024)
		n, _ := c.Read(bs)
		fmt.Println(i, string(bs[0:n]))
	}()

}
