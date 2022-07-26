package config

import (
	"encoding/json"
	"im-chatroom-broker/util"
	"os"
)

/**
{
  "port": "33121",
  "rocketmq": {
    "address" :"127.0.0.1:9876",
  },
  "redis": {
    "address":"47.95.148.121:6379",
    "password":"o1trUmeh",
    "db":1
  }
}
*/

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
	Port     string   `json:"port"`
	RocketMQ RocketMQ `json:"rocketmq"`
	Redis    Redis    `json:"redis"`
}

type RocketMQ struct {
	Address string `json:"address"`
}
type Redis struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	Db       int    `json:"db"`
}


var OP *Option


func NewDefaultOption() *Option {
	return &Option{
		Port: "33121",
		RocketMQ: RocketMQ{
			Address: "192.168.3.242:9876",
		},
		Redis: Redis{
			Address:  "47.95.148.121:6379",
			Password: "o1trUmeh",
			Db:       1,
		},
	}
}
