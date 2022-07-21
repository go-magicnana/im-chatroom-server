package main

import (
	"github.com/labstack/echo/middleware"
	"im-chatroom-gateway/domains"
	"im-chatroom-gateway/zaplog"
	"net/http"

	//导入echo包
	"github.com/labstack/echo"
	"im-chatroom-gateway/controllers"
)

func main() {
	//实例化echo对象。
	e := echo.New()
	e.Validator = &domains.CustomValidator{}

	e.Use(middleware.RequestID())
	//e.Use(lecho.Middleware(lecho.Config{
	//	Logger: logger,
	//}))

	// 定义post请求, url为：/getToken, 绑定getToken控制器函数
	e.POST("/imchatroom/token", controllers.GetToken)
	// 定义get请求
	e.GET("/imchatroom/config", controllers.GetConfig)

	e.POST("/imchatroom/room/create", controllers.CreateChatroom)

	e.POST("/imchatroom/message/push", controllers.MessagePush)

	//启动http server, 并监听1324端口，冒号（:）前面为空的意思就是绑定网卡所有Ip地址，本机支持的所有ip地址都可以访问。
	e.Logger.Fatal(e.Start(":33110"))
}

type handler func(echo.Context) (controllers.ApiResult, error)

func dispatcher(c echo.Context) error {

	url := c.Request().RequestURI
	var apiResult controllers.ApiResult
	parms, _ := c.FormParams()
	zaplog.Logger.Debugf("%s %v", url, parms)

	var ret any
	var e error

	switch url {
	case "/imchatroom/token":
		//ret, e = controllers.GetToken(c)
		break
	}

	zaplog.Logger.Infof("%s %v %v", url, ret, e)

	if e != nil {
		apiResult = controllers.NewApiResultError(e)
	} else {
		apiResult = controllers.NewApiResultOK(ret)
	}

	return c.JSON(http.StatusOK, apiResult)
}
