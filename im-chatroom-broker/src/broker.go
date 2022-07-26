package main

import (
	"im-chatroom-broker/config"
	"im-chatroom-broker/server"
	"im-chatroom-broker/zaplog"
)


func main(){
	zaplog.InitLogger()
	config.OP = loadConfig()
	server.Start()
}

func loadConfig() *config.Option{
	s := config.LoadConf("../conf/conf.json")
	zaplog.Logger.Infof("Init configuration %v",s)
	return s
}
