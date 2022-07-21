package util

import (
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	"testing"
)

func TestExternalIP(t *testing.T) {
	m := treemap.NewWith(utils.Int64Comparator)

	m.Put(12, "haha")
	m.Put(12, "hehe")
	m.Put(13, "gaga")

	m.Each(func(key interface{}, value interface{}) {
		fmt.Println(key, value)
	})
}
