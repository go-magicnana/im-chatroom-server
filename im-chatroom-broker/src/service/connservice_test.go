package service

import (
	"fmt"
	"testing"
)

func TestSetUserContext(t *testing.T) {
	SetRoomClient("1","clientName","userId")
	v := GetRoomMap("1")
	fmt.Println(v)
	v.Range(func(key, value any) bool {
		fmt.Println(key,value)
		return true
	})
}
