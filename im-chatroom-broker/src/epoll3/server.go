package epoll3

import "sync"

type Server struct {
	netpoll        netpoll      // 具体操作系统网络实现
	options        *options     // 服务参数
	readBufferPool *sync.Pool   // 读缓存区内存池
	handler        Handler      // 注册的处理
	ioEventQueues  []chan event // IO事件队列集合
	ioQueueNum     int32        // IO事件队列集合数量
	conns          sync.Map     // TCP长连接管理
	connsNum       int64        // 当前建立的长连接数量
	stop           chan int     // 服务器关闭信号
}
