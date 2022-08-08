package deliver

import (
	"im-chatroom-broker/ctx"
	"im-chatroom-broker/zaplog"
	"sync"
)

type pool struct {
	workers   int
	maxTasks  int
	taskQueue chan *ctx.Context

	mu     sync.Mutex
	closed bool
	done   chan struct{}
}

func newPool(w int, t int) *pool {
	return &pool{
		workers:   w,
		maxTasks:  t,
		taskQueue: make(chan *ctx.Context, t),
		done:      make(chan struct{}),
	}
}

func (p *pool) Close() {
	p.mu.Lock()
	p.closed = true
	close(p.done)
	close(p.taskQueue)
	p.mu.Unlock()
}

func (p *pool) addTask(conn *ctx.Context) {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return
	}
	p.mu.Unlock()

	p.taskQueue <- conn
}

func (p *pool) start() {
	for i := 0; i < p.workers; i++ {
		go p.startWorker()
	}
}

func (p *pool) startWorker() {
	zaplog.Logger.Infof("Worker start")

	for {
		select {
		case <-p.done:
			return
		case conn := <-p.taskQueue:
			if conn != nil {
				//buf, ok := conn.Queue.Dequeue()
				//if ok {
				//	b := buf.([]byte)
				//	conn.Conn.AsyncWrite(b, nil)
				//}
			}
		}
	}
}
