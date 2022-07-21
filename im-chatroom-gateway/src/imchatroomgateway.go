package main

import (
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/ziflex/lecho/v3"
	"im-chatroom-gateway/domains"
	"os"

	//导入echo包
	"github.com/labstack/echo"
	"im-chatroom-gateway/controllers"
)

func main() {
	boot.Start()
}
