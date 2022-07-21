package controllers

import (
	"context"
	"github.com/labstack/echo"
	"github.com/ziflex/lecho/v3"
	"im-chatroom-gateway/redis"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	BrokerInfo     string = "imchatroom:brokerinfo:"
	BrokerInstance string = "imchatroom:brokerinstance"
	BrokerCapacity string = "imchatroom:brokercapacity:"
)

func GetConfig(ct echo.Context) error {
	e := echo.New()
	e.Logger = lecho.New(
		os.Stdout,
		lecho.WithFields(map[string]interface{}{"name": "lecho factory"}),
		lecho.WithTimestamp(),
		lecho.WithCaller(),
		lecho.WithPrefix("controllers.GetConfig"),
	)

	// 先获取所有服务器列表
	serverlist, err := redis.Rdb.SMembers(context.Background(), BrokerInstance).Result()

	// 遍历列表，获取每台服务加入量
	slist := []string{}
	if err != redis.Nil && len(serverlist) > 0 {

		for i := range serverlist {
			sip := serverlist[i]
			serverCap, err := redis.Rdb.SCard(context.Background(), BrokerCapacity+sip).Result()
			if err != redis.Nil {
				slist = append(slist, strconv.FormatInt(serverCap, 10)+"-"+sip)
			} else {
				slist = append(slist, "0-"+sip)
			}
		}
		sort.Strings(slist)

	}
	e.Logger.Info("get serverBrokers :", slist)
	var appConfig AppConfig
	var servers []serverInfo
	for i := range slist {
		svr := slist[i]
		svrArr := strings.Split(svr[strings.Index(svr, "-")+1:], ":")
		svrTmp := serverInfo{
			Ip:   svrArr[0],
			Port: svrArr[1],
		}
		servers = append(servers, svrTmp)
	}
	appConfig.Servers = servers

	appConfig.HeartTime = 5
	return ct.JSON(http.StatusOK, NewApiResultOK(appConfig))
	// return NewApiResultOK(appConfig)
}

type AppConfig struct {
	Servers   []serverInfo `json:"servers"`
	HeartTime int8         `json:"heartTime"`
}

type serverInfo struct {
	Ip   string `json:"ip"`
	Port string `json:"port"`
}
