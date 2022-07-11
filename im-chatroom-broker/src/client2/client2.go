package client2

import (
	"fmt"
	"im-chatroom-broker/protocol"
	"net"
	"os"
)

func Start() {
	server := "localhost:33121"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	fmt.Println("connect success")
	send(conn)
}

func send(conn net.Conn) {
	p := protocol.Packet{
		MessageId: "e10adc3949ba59abbe56e057f20f883e",
		Version:   1,
		Command:   1000,
	}

}
