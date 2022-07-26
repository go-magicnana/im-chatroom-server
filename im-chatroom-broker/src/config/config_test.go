package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"testing"
)

func TestLoadConf(t *testing.T) {
	s := LoadConf("../../conf/conf-test.json")
	fmt.Println(s)
}

func TestViper(t *testing.T) {
	fmt.Println(os.Getenv("GO_ENV"))

	endpoint := viper.GetString("rocketmq.endpoint")

	fmt.Println(endpoint)

	//fmt.Println(profile.GetStringValue("rocketmq.endpoint"))
}
