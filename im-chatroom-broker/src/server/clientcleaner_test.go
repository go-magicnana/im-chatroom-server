package server

import (
	"sync"
	"testing"
)

func TestClientCleanerTask(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	ClientCleanerTask()
	wg.Wait()
}
