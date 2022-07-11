package util

import (
	"fmt"
	"testing"
)

func TestExternalIP(t *testing.T) {
	ip,_:=ExternalIP()
	fmt.Println(ip.String())

}

