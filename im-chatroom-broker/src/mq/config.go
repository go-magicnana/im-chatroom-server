package mq

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func init() {

	os.Setenv("GO_ENV", "dev")

	viper.SetConfigFile("../conf/config-" + os.Getenv("GO_ENV") + ".json")
	viper.AutomaticEnv()
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
	}

	fmt.Println(viper.GetViper())

	//if err := viper.Unmarshal(&config); err != nil {
	//return nil, err
	//}

}
