//go:build linux
// +build linux

package epoll

import (
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"im-chatroom-broker/zaplog"
	"strconv"
	"strings"
	"syscall"
)

// 对端关闭连接 8193
const (
	EpollRead  = syscall.EPOLLIN | syscall.EPOLLPRI | syscall.EPOLLERR | syscall.EPOLLHUP | unix.EPOLLET | syscall.EPOLLRDHUP
	EpollClose = uint32(syscall.EPOLLIN | syscall.EPOLLRDHUP)
)


type epoll struct {
	listenFD int
	epollFD  int
}

func newNetpoll(address string) (netpoll, error) {
	listenFD, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		zaplog.Logger.Errorf("NewEpoll Init error %v",err)
		return nil, err
	}

	zaplog.Logger.Infof("NewEpoll Init OK %v",listenFD)

	err = syscall.SetsockoptInt(listenFD, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		return nil, err
	}

	addr, port, err := getIPPort(address)
	if err != nil {
		return nil, err
	}
	err = syscall.Bind(listenFD, &syscall.SockaddrInet4{
		Port: port,
		Addr: addr,
	})
	if err != nil {
		zaplog.Logger.Errorf("NewEpoll Bind Error %v %v",listenFD,err)
		return nil, err
	}
	zaplog.Logger.Infof("NewEpoll Bind OK %v",listenFD)



	err = syscall.Listen(listenFD, 1024)
	if err != nil {
		zaplog.Logger.Errorf("NewEpoll Listen Error %v %v",listenFD,err)
		return nil, err
	}

	zaplog.Logger.Infof("NewEpoll Listen OK %v",listenFD)


	epollFD, err := syscall.EpollCreate1(0)
	if err != nil {
		zaplog.Logger.Errorf("NewEpoll EpollCreate1 Error %v %v",epollFD,err)
		return nil, err
	}

	zaplog.Logger.Infof("NewEpoll EpollCreate1 OK %v %v",epollFD,0)


	return &epoll{listenFD: listenFD, epollFD: epollFD}, nil
}

func (n *epoll) accept() (nfd int, addr string, err error) {
	nfd, sa, err := syscall.Accept(n.listenFD)
	if err != nil {
		zaplog.Logger.Errorf("NewEpoll Accept Error %v %v",n.listenFD,err)
		return
	}

	zaplog.Logger.Infof("NewEpoll Accept OK %v %v %v",n.listenFD,nfd,sa)


	// 设置为非阻塞状态
	err = syscall.SetNonblock(nfd, true)
	if err != nil {
		return
	}

	err = syscall.EpollCtl(n.epollFD, syscall.EPOLL_CTL_ADD, nfd, &syscall.EpollEvent{
		Events: EpollRead,
		Fd:     int32(nfd),
	})
	if err != nil {
		return
	}

	zaplog.Logger.Infof("NewEpoll AddReadEvent OK %v %v",n.epollFD,nfd)


	s := sa.(*syscall.SockaddrInet4)
	addr = fmt.Sprintf("%d.%d.%d.%d:%d", s.Addr[0], s.Addr[1], s.Addr[2], s.Addr[3], s.Port)
	return
}

func (n *epoll) addRead(fd int) error {
	err := syscall.EpollCtl(n.epollFD, syscall.EPOLL_CTL_ADD, fd, &syscall.EpollEvent{
		Events: EpollRead,
		Fd:     int32(fd),
	})
	if err != nil {
		return err
	}
	return nil
}

func (n *epoll) closeFD(fd int) error {
	// 移除文件描述符的监听
	err := syscall.EpollCtl(n.epollFD, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		return err
	}

	// 关闭文件描述符
	err = syscall.Close(fd)
	if err != nil {
		return err
	}

	return nil
}

func (n *epoll) getEvents() ([]event, error) {
	zaplog.Logger.Debugf("NewEpoll EpollWait start")

	epollEvents := make([]syscall.EpollEvent, 100)
	num, err := syscall.EpollWait(n.epollFD, epollEvents, -1)
	zaplog.Logger.Errorf("NewEpoll EpollWait OK %v %v",num,err)

	if err != nil {
		return nil, err
	}

	events := make([]event, 0, len(epollEvents))
	for i := 0; i < num; i++ {
		event := event{
			FD: epollEvents[i].Fd,
		}
		if epollEvents[i].Events == EpollClose {
			event.Type = EventClose
		} else {
			event.Type = EventIn
		}
		events = append(events, event)
	}

	return events, nil
}

func (n *epoll) closeFDRead(fd int) error {
	_, _, e := syscall.Syscall(syscall.SHUT_RD, uintptr(fd), 0, 0)
	if e != 0 {
		return e
	}
	return nil
}

func getIPPort(addr string) (ip [4]byte, port int, err error) {
	strs := strings.Split(addr, ":")
	if len(strs) != 2 {
		err = errors.New("addr error")
		return
	}

	if len(strs[0]) != 0 {
		ips := strings.Split(strs[0], ".")
		if len(ips) != 4 {
			err = errors.New("addr error")
			return
		}
		for i := range ips {
			data, err := strconv.Atoi(ips[i])
			if err != nil {
				return ip, 0, err
			}
			ip[i] = byte(data)
		}
	}

	port, err = strconv.Atoi(strs[1])
	return
}
