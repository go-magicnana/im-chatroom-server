package boot

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"im-chatroom-gateway/controllers"
	"im-chatroom-gateway/domains"
)

func Start() {
	//实例化echo对象。
	e := echo.New()
	e.Validator = &domains.CustomValidator{}

	e.Use(middleware.RequestID())
	//e.Use(lecho.Middleware(lecho.Config{
	//	Logger: logger,
	//}))

	/**
		TypeNoticeBlockUser   = 3103
	TypeNoticeUnblockUser = 3104
	TypeNoticeCloseRoom   = 3105
	TypeNoticeBlockRoom   = 3106
	TypeNoticeUnblockRoom = 3107
	*/

	// 定义post请求, url为：/getToken, 绑定getToken控制器函数
	e.POST("/imchatroom/token", controllers.GetToken)
	// 定义get请求
	e.GET("/imchatroom/config", controllers.GetConfig)

	e.POST("/imchatroom/room/create", controllers.CreateChatroom)

	e.POST("/imchatroom/notice/block_user", controllers.NoticeBlockUser)
	e.POST("/imchatroom/notice/unblock_user", controllers.NoticeUnblockUser)
	e.POST("/imchatroom/notice/close_room", controllers.NoticeCloseRoom)
	e.POST("/imchatroom/notice/block_room", controllers.NoticeBlockRoom)
	e.POST("/imchatroom/notice/unblock_room", controllers.NoticeUnblockRoom)

	//启动http server, 并监听1324端口，冒号（:）前面为空的意思就是绑定网卡所有Ip地址，本机支持的所有ip地址都可以访问。
	e.Start(":33110")
}
