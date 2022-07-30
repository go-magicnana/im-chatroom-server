package handler

import (
	"fmt"
	"testing"
)

func Test_spiltList(t *testing.T) {

	s1 := make([]interface{},0)
	s1 = append(s1,"haha1")
	s1 = append(s1,"haha2")
	s1 = append(s1,"haha3")
	s1 = append(s1,"haha4")
	s1 = append(s1,"haha5")

	for i := 1; i < 10 ; i++ {
		fmt.Println(spiltList(s1,i))
	}
}
