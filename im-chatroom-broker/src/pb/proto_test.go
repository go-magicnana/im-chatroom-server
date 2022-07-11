package pb

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"testing"
)

func TestUser(t *testing.T) {
	user := User{
		UserId: 1,
		Username: "testname",
	}
	// 将protocol消息序列化成二进制
	userByte, err := proto.Marshal(&user)
	if nil != err {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(userByte)

	// 将二进制消息反序列化
	userNew := User{}
	err2 := proto.Unmarshal(userByte, &userNew)
	if nil != err2 {
		fmt.Println(err2.Error())
	}

	fmt.Printf("%+v\n", userNew)
}
