package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo"
	"github.com/ziflex/lecho/v3"
	"net/http"
	"os"
	"time"
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
		ct.JSON(http.StatusOK, gin.H{"code": 1001, "message": "param is error"})
	}

	// 创建roomId
	timeUnix := time.Now().UnixNano() / 1e6
	roomId := userId + fmt.Sprintf("%d", timeUnix)

	return ct.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": roomId})
}
