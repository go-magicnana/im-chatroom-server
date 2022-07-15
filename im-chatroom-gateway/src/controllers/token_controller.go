package controllers

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo"
	"github.com/ziflex/lecho/v3"
	"im-chatroom-gateway/src/rediss"
	"net/http"
	"os"
	"strconv"
	"time"
)

const userTokenKey string = "imchatroom:userauth:token:"

func GetToken(ct echo.Context) error {
	e := echo.New()
	e.Logger = lecho.New(
		os.Stdout,
		lecho.WithFields(map[string]interface{}{"name": "lecho factory"}),
		lecho.WithTimestamp(),
		lecho.WithCaller(),
		lecho.WithPrefix("controllers.GetToken"),
	)

	//获取post请求的表单参数
	userId := ct.FormValue("userId")
	name := ct.FormValue("name")
	avatar := ct.FormValue("avatar")
	sexstr := ct.FormValue("sex")

	if userId == "" {
		e.Logger.Info("userId is illegal")
		ct.JSON(http.StatusOK, gin.H{"code": 1001, "message": "param is error"})
	}

	userinfo := UserInfo{}
	userinfo.UserId = userId
	userinfo.Name = name
	userinfo.Avatar = avatar
	sex, _ := strconv.Atoi(sexstr)
	userinfo.Sex = sex

	timeUnix := time.Now().Unix()
	// userToken:userId时间戳 获取md5值作为token
	data := []byte("userToken:" + userId + fmt.Sprintf("%d", timeUnix))
	// 将[]byte转成16进制
	userToken := fmt.Sprintf("%x", md5.Sum(data))

	// 存入redis 需要序列化
	data, err := json.Marshal(userinfo)
	if err != nil {
		e.Logger.Error("json userinfo occur err")
		return ct.JSON(http.StatusOK, gin.H{"code": 1001, "message": "Server Error"})
	}
	result := rediss.RedisSingleton().Set(context.Background(), userTokenKey+userToken, data, time.Hour*24)
	e.Logger.Info(result)
	return ct.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": userToken})
}

type UserInfo struct {
	UserId string `json:"userId"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Sex    int    `json:"sex"`
}
