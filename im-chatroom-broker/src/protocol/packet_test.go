package protocol

import (
	"fmt"
	"testing"
)

func TestNewResponseOK(t *testing.T) {

	header := MessageHeader{
		MessageId: "e10adc3949ba59abbe56e057f20f883a",
		Command:   CommandSignal,
		Flow:      FlowUp,
		Type:      TypeSignalLogin,
	}

	body := MessageBodySignalLogin{
		Token: "token1",
	}

	packet := Packet{
		Header: header, Body: body,
	}


	p := NewResponseOK(&packet,nil)
	fmt.Println(p.ToString())
}
