package zaplog

import (
	"fmt"
	"testing"
)

func TestTest1(t *testing.T) {
	InitLogger()
	defer Logger.Sync()
	Infof("这是一个描述 %s","binggo")

}

func TestFoo(t *testing.T)  {
	fmt.Println("haha")
}
