package config

import (
	"fmt"
	"testing"
)

func TestLoadConf(t *testing.T) {
	s := LoadConf("../../conf/conf.json")
	fmt.Println(s)
}
