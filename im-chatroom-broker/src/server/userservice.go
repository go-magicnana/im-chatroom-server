package server

import (
	"context"
	"fmt"
	context2 "im-chatroom-broker/context"
	"im-chatroom-broker/redis"
	"im-chatroom-broker/util"
	"sync"
	"time"
)

var users sync.Map
var dirty sync.Map

const (
	UserHash  string = "dudu:broker:user_hash"
	UserOther string = "dudu:broker:user_other"
)
//
//func setDirtyConnection(c *context2.Context){
//	dirty.Store(c.RemoteAddr,c)
//}

func SetUser(userId string, c *context2.Context) {

	users.Store(userId, c)

	//cmd1 := redis.Rdb.HSet(c.Ctx, UserHash, userId, c.Broker)
	//if cmd1.Err() != nil {
	//	Panic(cmd1.Err())
	//}


}

func GetUserLocal(userId string) (*context2.Context, bool) {
	ay, exist := users.Load(userId)
	if exist {
		return ay.(*context2.Context), true
	} else {
		return nil, false
	}
}

func UserLocal2String(){

	go func() {
		for {
			fmt.Print("当前 ")
			users.Range(func(key, value any) bool {
				fmt.Print(key)
				fmt.Print("\t")
				return true
			})
			fmt.Println("")
			time.Sleep(time.Second)
		}
		//wg.Done()
	}()

}

func GetUserGlobal(ctx context.Context, userId string) (string, bool) {

	ret := redis.Rdb.HGet(ctx, UserHash, userId)

	if ret.Err() != nil {
		util.Panic(ret.Err())
	}

	if ret != nil {
		return ret.Val(), true
	}
	return "", false
}

//TODO 单一用户多客户端登录
