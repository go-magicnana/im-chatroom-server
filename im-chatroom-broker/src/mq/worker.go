package mq

import (
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/zaplog"
	"sync"
)

type pool struct {
	workers   int
	maxTasks  int
	taskQueue chan *protocol.PacketDeliver

	mu     sync.Mutex
	closed bool
	done   chan struct{}
}

func newPool(w int, t int) *pool {
	return &pool{
		workers:   w,
		maxTasks:  t,
		taskQueue: make(chan *protocol.PacketDeliver, t),
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

func (p *pool) addTask(packet *protocol.PacketDeliver) {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return
	}

	p.taskQueue <- packet
	p.mu.Unlock()
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
		case pp := <-p.taskQueue:
			if pp != nil {

				deliver(pp.ToMQ,pp.Packet)

			}
		}
	}
}
