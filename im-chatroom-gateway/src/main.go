package main

import (
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/ziflex/lecho/v3"
	"os"

	//导入echo包
	"github.com/labstack/echo"
	"im-chatroom-gateway/src/controllers"
)

func main() {
	//实例化echo对象。
	e := echo.New()

	logger := lecho.New(
		os.Stdout,
		lecho.WithLevel(log.INFO),
		lecho.WithTimestamp(),
		lecho.WithCaller(),
	)
	e.Logger = logger

	e.Use(middleware.RequestID())
	//e.Use(lecho.Middleware(lecho.Config{
	//	Logger: logger,
	//}))

	// 定义post请求, url为：/getToken, 绑定getToken控制器函数
	e.POST("/getToken", controllers.GetToken)
	// 定义get请求
	e.GET("/config", controllers.GetConfig)

	e.POST("/room/create", controllers.CreateChatroom)

	//启动http server, 并监听1324端口，冒号（:）前面为空的意思就是绑定网卡所有Ip地址，本机支持的所有ip地址都可以访问。
	e.Logger.Fatal(e.Start(":1324"))
}
