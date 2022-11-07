package config

import (
	"fmt"
	"im-chatroom-broker/protocol"
	"testing"
)

func TestLoadConf(t *testing.T) {
	header := protocol.MessageHeader{
		MessageId: "LoginMessageId-login",
		Command:   protocol.CommandSignal,
		Flow:      protocol.FlowUp,
		Type:      protocol.TypeSignalLogin,
	}

	body := protocol.MessageBodySignalLogin{
		Token:  "token",
		Device: "MAC",
		//UserId: "1001",
		//Name:   "张三丰",
		//Avatar: "https://img1.baidu.com/it/u=2848117662,2869906655&fm=253&fmt=auto&app=138&f=JPEG?w=500&h=501",
	}

	packet := protocol.Packet{
		Header: header, Body: body,
	}

	fmt.Println(packet.ToString())
}
