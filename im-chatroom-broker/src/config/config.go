package config

import (
	"encoding/json"
	"im-chatroom-broker/util"
	"os"
)

func LoadConf(path string) *Option {

	file, e := os.Open(path)
	defer file.Close()

	if e != nil {
		util.Panic(e)
	}

	option := Option{}
	if err := json.NewDecoder(file).Decode(&option); err != nil {
		util.Panic(err)
	}

	return &option
}

type Option struct {
	RocketMQ string `json:"rocketmq"`
	Obj Obj `json:"obj"`
}

type Obj struct {
	Name string `json:"name"`
}
