//go:build !linux
// +build !linux

package epoll

import "im-chatroom-broker/util"

func newNetpoll(address string) (netpoll, error) {
	util.Panic("please run on linux")
	return nil,nil
}
