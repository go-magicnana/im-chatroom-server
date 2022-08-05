package main

import "im-chatroom-broker/core"

//func main() {
//	epoll2.StartServer()
//}

//func main() {
//	var handler = &netpoll.DataHandler{
//		NoShared:   true,
//		NoCopy:     true,
//		BufferSize: 1024,
//		HandlerFunc: func(req []byte) (res []byte) {
//			res = req
//			return
//		},
//	}
//	if err := netpoll.ListenAndServe("tcp", ":33121", handler); err != nil {
//		util.Panic(err)
//	}
//}

func main() {
	core.Start()
}
