package config

import (
	"encoding/json"
	"im-chatroom-gateway/util"
	"os"
	"sync"
)

var once sync.Once

var OP *Option

func init() {

	os.Setenv("ROCKETMQ_GO_LOG_LEVEL", "error")

	OP = Singleton()
}

func Singleton() *Option {
	once.Do(func() {
		OP = LoadConf("../conf/conf.json")

	})

	return OP
}
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
