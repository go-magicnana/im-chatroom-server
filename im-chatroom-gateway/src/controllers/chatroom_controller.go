package controllers

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/ziflex/lecho/v3"
	"golang.org/x/net/context"
	"im-chatroom-gateway/apierror"
	"im-chatroom-gateway/redis"
	"net/http"
	"os"
	"time"
)

const (
	RoomInfo string = "imchatroom:roominfo:"
)

func CreateChatroom(ct echo.Context) error {
	e := echo.New()
	e.Logger = lecho.New(
		os.Stdout,
		lecho.WithFields(map[string]interface{}{"name": "lecho factory"}),
		lecho.WithTimestamp(),
		lecho.WithCaller(),
		lecho.WithPrefix("controllers.CreateChatroom"),
	)

	//获取post请求的表单参数
	userId := ct.FormValue("userId")
	if userId == "" {
		e.Logger.Info("userId is illegal")
		return ct.JSON(http.StatusOK, NewApiResultError(apierror.InvalidParameter))
	}

	// 创建roomId
	timeUnix := time.Now().UnixNano() / 1e6
	roomId := userId + fmt.Sprintf("%d", timeUnix)

	return ct.JSON(http.StatusOK, NewApiResultOK(roomId))
}

func GetRoomMembers(ct echo.Context) error {
	e := echo.New()
	e.Logger = lecho.New(
		os.Stdout,
		lecho.WithFields(map[string]interface{}{"name": "lecho factory"}),
		lecho.WithTimestamp(),
		lecho.WithCaller(),
		lecho.WithPrefix("controllers.CreateChatroom"),
	)

	//获取post请求的表单参数
	roomId := ct.FormValue("roomId")
	if roomId == "" {
		e.Logger.Info("roomId is illegal")
		return ct.JSON(http.StatusOK, NewApiResultError(apierror.InvalidParameter))
	}

	snum, err := redis.Rdb.SCard(context.Background(), RoomInfo+roomId).Result()
	if err == redis.Nil {
		snum = 0
	}

	return ct.JSON(http.StatusOK, NewApiResultOK(snum))
}
