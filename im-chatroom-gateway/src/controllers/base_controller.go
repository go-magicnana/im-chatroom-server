package controllers

import (
	"context"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"github.com/labstack/echo"
	"im-chatroom-gateway/apierror"
	"im-chatroom-gateway/redis"
	"im-chatroom-gateway/zaplog"
	"net/http"
)

const (
	BrokerInstance string = "imchatroom:brokerinstance"
	BrokerCapacity string = "imchatroom:brokercapacity:"
)

func GetConfig(c echo.Context) error {

	a, _ := c.FormParams()
	zaplog.Logger.Debugf("%s %v", c.Request().RequestURI, a)

	// 先获取所有服务器列表
	serverlist, err := redis.Rdb.SMembers(context.Background(), BrokerInstance).Result()
	if err != nil {
		return write(c, http.StatusOK, NewApiResultError(apierror.StorageResponseError.Format(err.Error())))
	}

	if serverlist==nil || len(serverlist)==0 {
		return write(c, http.StatusOK, NewApiResultError(apierror.StorageResponseEmpty))

	}

	var l []string

	m := treemap.NewWith(utils.Int64Comparator)


	// 遍历列表，获取每台服务加入量
		for _,v := range serverlist {
			cap, err := redis.Rdb.SCard(context.Background(), BrokerCapacity+v).Result()
			if err != nil {
				l = append(l, v)
			} else {
				m.Put(v,cap)
			}
		}



	//e.Logger.Info("get serverBrokers :", slist)
	//var appConfig AppConfig
	//var servers []serverInfo
	//for i := range slist {
	//	svr := slist[i]
	//	svrArr := strings.Split(svr[strings.Index(svr, "-")+1:], ":")
	//	svrTmp := serverInfo{
	//		Ip:   svrArr[0],
	//		Port: svrArr[1],
	//	}
	//	servers = append(servers, svrTmp)
	//}
	//appConfig.Servers = servers
	//
	//appConfig.HeartTime = 1
	//
	//return ct.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": appConfig})
	return nil
}

type AppConfig struct {
	Servers   []serverInfo `json:"servers"`
	HeartTime int8         `json:"heartTime"`
}

type serverInfo struct {
	Ip   string `json:"ip"`
	Port string `json:"port"`
}
