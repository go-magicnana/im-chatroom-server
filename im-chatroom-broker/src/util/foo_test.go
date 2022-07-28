package util

import (
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"golang.org/x/net/context"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func Test1(t *testing.T) {

	fmt.Println(strconv.FormatInt(999999999999999999, 10))
	fmt.Println(strconv.FormatInt(11111111111111111, 10))
	fmt.Println(strconv.FormatInt(1111111111111111, 10))
	fmt.Println(strconv.FormatInt(111111111111111, 10))
	fmt.Println(strconv.FormatInt(11111111111111, 10))
	fmt.Println(strconv.FormatInt(1111111111111, 10))
	fmt.Println(strconv.FormatInt(111111111111, 10))
	fmt.Println(strconv.FormatInt(11111111111, 10))
	fmt.Println(strconv.FormatInt(1111111111, 10))
	fmt.Println(strconv.FormatInt(111111111, 10))
	fmt.Println(strconv.FormatInt(11111111, 10))
	fmt.Println(strconv.FormatInt(1111111, 10))
	fmt.Println(strconv.FormatInt(111111, 10))
	fmt.Println(strconv.FormatInt(11111, 10))
	fmt.Println(strconv.FormatInt(1111, 10))
	fmt.Println(strconv.FormatInt(111, 10))
	fmt.Println(strconv.FormatInt(11, 10))
	fmt.Println(strconv.FormatInt(1, 10))
}

func TestStartWith(t *testing.T) {
	s := "auth sdfsdfsdf"

	fmt.Println(strings.Index(s, "auth"))
	fmt.Println(s[0:4])
	fmt.Println(s[4:])
}

func TestMillion(t *testing.T) {
	fmt.Println(time.Now().UnixNano() / 1e6)
}

func TestSyncMap(t *testing.T) {
	var users sync.Map
	fmt.Println(users)

	var wg sync.WaitGroup
	wg.Add(10)

	go func() {
		value := "你好张三"
		users.Store("zhangsan", &value)
		users.Store("wangwu", &value)

		wg.Done()
	}()

	go func() {
		value := "你好李四"
		users.Store("lisi", &value)
		wg.Done()
	}()

	go func() {
		for {
			fmt.Println(users)
			time.Sleep(time.Second)
		}
		wg.Done()
	}()

	go func() {
		for {
			users.Load("zhangsan")
		}
	}()

	go func() {
		for {
			users.Load("lisi")
		}
	}()
	wg.Wait()
	fmt.Println("OK")
}

func TestMd5Length(t *testing.T) {
	s := "e10adc3949ba59abbe56e057f20f883e"
	fmt.Println(len(S2b(s)))
}

func TestSlice(t *testing.T) {

	s := "e10adc3949ba59abbe56e057f20f883e"

	bs := S2b(s)
	fmt.Println(len(bs), cap(bs))
	var b byte = 1
	bs = append(bs, b)
	fmt.Println(len(bs), cap(bs))
	//fmt.Println(B2s(bs))

	s1 := bs[0:32]
	s2 := bs[32:]
	fmt.Println(B2s(s1))
	fmt.Println(s2)

}

func TestCreateTopic(t *testing.T) {

	topicName := "imchatroom_one_topic_192_168_1_73_33121"

	endPoint := []string{"47.95.149.47:9876"}
	// 创建主题
	testAdmin, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver(endPoint)))
	if err != nil {
		fmt.Printf("connection error: %s\n", err.Error())
	}
	err = testAdmin.CreateTopic(context.Background(), admin.WithTopicCreate(topicName),admin.WithBrokerAddrCreate("47.95.149.47:9876"))
	if err != nil {
		fmt.Printf("createTopic error: %s\n", err.Error())
	}
}
