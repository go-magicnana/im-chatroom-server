package controllers

import (
	"context"
	"github.com/labstack/echo"
	"im-chatroom-gateway/apierror"
	"im-chatroom-gateway/redis"
	"im-chatroom-gateway/zaplog"
	"net/http"
	"sort"
	"strings"
)

const (
	BrokerInfo     string = "imchatroom:broker.info:"
	BrokerInstance string = "imchatroom:broker.instance"
	BrokerCapacity string = "imchatroom:broker.capacity:"
)

func GetConfig(c echo.Context) error {

	//a, _ := c.FormParams()
	zaplog.Logger.Info("request url:", c.Request().RequestURI, "...")

	// 先获取所有服务器列表
	serverList, err := redis.Rdb.SMembers(context.Background(), BrokerInstance).Result()
	if err != nil {
		return write(c, http.StatusOK, NewApiResultError(apierror.StorageResponseError.Format(err.Error())))
	}

	if serverList == nil || len(serverList) == 0 {
		return write(c, http.StatusOK, NewApiResultError(apierror.StorageResponseEmpty))
	}

	var brokerList []BrokerCompare

	// 遍历列表，获取每台服务加入量
	for _, v := range serverList {
		server := strings.Split(v, ":")
		cap, err := redis.Rdb.SCard(context.Background(), BrokerCapacity+v).Result()
		if err != nil {
			brokerList = append(brokerList, BrokerCompare{999999999, server[0], server[1]})
		} else {
			brokerList = append(brokerList, BrokerCompare{cap, server[0], server[1]})
		}
	}

	sortBySize(brokerList)
	config := AppConfig{
		Servers:   brokerList,
		HeartTime: 10,
	}

	return write(c, http.StatusOK, NewApiResultOK(config))

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
}

type AppConfig struct {
	Servers   []BrokerCompare `json:"servers"`
	HeartTime int32           `json:"heartTime"`
}

type serverInfo struct {
	Ip   string `json:"ip"`
	Port string `json:"port"`
}

//通用排序
//结构体排序，必须重写数组Len() Swap() Less()函数
type body_wrapper struct {
	Bodys []BrokerCompare
	by    func(p, q *BrokerCompare) bool //内部Less()函数会用到
}
type SortBodyBy func(p, q *BrokerCompare) bool //定义一个函数类型

//数组长度Len()
func (acw body_wrapper) Len() int {
	return len(acw.Bodys)
}

//元素交换
func (acw body_wrapper) Swap(i, j int) {
	acw.Bodys[i], acw.Bodys[j] = acw.Bodys[j], acw.Bodys[i]
}

//比较函数，使用外部传入的by比较函数
func (acw body_wrapper) Less(i, j int) bool {
	return acw.by(&acw.Bodys[i], &acw.Bodys[j])
}

//自定义排序字段，参考SortBodyByCreateTime中的传入函数
func SortBody(bodys []BrokerCompare, by SortBodyBy) {
	sort.Sort(body_wrapper{bodys, by})
}

//按照createtime排序，需要注意是否有createtime
func sortBySize(bodys []BrokerCompare) {
	sort.Sort(body_wrapper{bodys, func(p, q *BrokerCompare) bool {
		return p.Size < q.Size
	}})
}

type BrokerCompare struct {
	Size int64  `json:"size"`
	Ip   string `json:"ip"`
	Port string `json:"port"`
}
