package ctx

import (
	"fmt"
	"sync"
	"testing"
)

func TestContext_ToString(t *testing.T) {

	wg := new(sync.WaitGroup)
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go createContext(wg, i)
	}



	wg.Wait()
	wg.Add(1)
	go queryContext(wg)
	wg.Wait()

}

func queryContext(wg *sync.WaitGroup) {
	for i := 0; i < 10; i++ {
		clientName := "client-" + fmt.Sprintf("%d", i)
		cc := GetContext(clientName)
		if cc == nil {
			fmt.Println(clientName+" not exist")
		}else{
			fmt.Println(cc.Broker)
		}
	}

	wg.Done()
}

func createContext(wg *sync.WaitGroup, i int) {

	clientName := "client-" + fmt.Sprintf("%d", i)

	cc := &Context{
		ClientName: clientName,
		Broker:     clientName,
	}

	fmt.Println(cc)

	wg.Done()
}

