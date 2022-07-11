package server

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var zhangsan *Context = NewContext("", "borkerAddr", nil, nil, nil)
var lisi *Context = NewContext("", "borkerAddr", nil, nil, nil)

func TestUserService(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		SetUser("zhangsan", zhangsan)
	}()

	go func() {
		SetUser("lisi", lisi)
	}()

	go func() {
		time.Sleep(time.Second * 10)
		SetUser("wangwu", zhangsan)
	}()
	UserLocal2String()

	wg.Wait()
	fmt.Println("OK")
}
