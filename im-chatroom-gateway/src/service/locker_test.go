package service

import (
	"fmt"
	"golang.org/x/net/context"
	"testing"
)

func TestLock(t *testing.T) {
	ret := Lock(context.Background(),"1212")
	fmt.Println(ret)
}
