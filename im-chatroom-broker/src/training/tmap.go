package training

import (
	"fmt"
	"sync"
)

var m sync.Map

func go1(key string,value chan string)  {



	m.Store(key,value)

}

func go2(key string)  {
	m.Delete(key)
}


func go3(){
	m.Range(func(key, value any) bool {
		fmt.Printf("%s %p \n",key,value)
		return true
	})
}

