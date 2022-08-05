package epoll2

import (
	"flag"
	"im-chatroom-broker/util"
	"log"
	"net"
	_ "net/http/pprof"
)

var (
	c = flag.Int("c", 10, "concurrency")
)

var epoller *epoll
var workerPool *pool

func StartServer() {
	flag.Parse()

	ln, err := net.Listen("tcp", ":33121")
	if err != nil {
		util.Panic(err)
	}

	workerPool = newPool(*c, 8)
	workerPool.start()

	epoller, err = MkEpoll()
	if err != nil {
		util.Panic(err)
	}

	go start()

	for {
		conn, e := ln.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				log.Printf("accept temp err: %v", ne)
				continue
			}

			log.Printf("accept err: %v", e)
			return
		}

		if err := epoller.Add(conn); err != nil {
			log.Printf("failed to add connection %v", err)
			conn.Close()
		}
	}

	workerPool.Close()
}

func start() {
	for {
		connections, err := epoller.Wait()
		if err != nil {
			log.Printf("failed to epoll wait %v", err)
			continue
		}
		for _, conn := range connections {
			if conn == nil {
				break
			}

			workerPool.addTask(conn)
		}
	}
}
