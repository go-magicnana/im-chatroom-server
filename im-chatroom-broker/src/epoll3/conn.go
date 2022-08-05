package epoll3

import "sync"

type Conn struct {
	Server     *Server  // 服务器引用
	fd         int32    // 文件描述符
	RemoteAddr string   // 对端地址
	Data       sync.Map // 业务自定义数据，用作扩展
}
