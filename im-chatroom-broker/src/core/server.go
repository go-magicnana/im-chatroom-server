package core

import (
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/robfig/cron/v3"
	"im-chatroom-broker/config"
	"im-chatroom-broker/ctx"
	err "im-chatroom-broker/error"
	"im-chatroom-broker/handler"
	"im-chatroom-broker/protocol"
	"im-chatroom-broker/serializer"
	"im-chatroom-broker/service"
	"im-chatroom-broker/util"
	"im-chatroom-broker/zaplog"
)

type server struct {
	gnet.BuiltinEventEngine
	serializer *serializer.JsonSerializer
}

func (s *server) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {

	cc := &ctx.Context{
		Conn:       c,
		Broker:     ctx.BrokerAddress,
		ClientName: c.RemoteAddr().String(),
		Time:       util.CurrentSecond(),
	}

	c.SetContext(cc)
	ctx.OpenContext(c.RemoteAddr().String(), cc)

	zaplog.Logger.Debugf("Connection open %s", c.RemoteAddr().String())

	return
}

func (s *server) OnClose(c gnet.Conn, err error) (action gnet.Action) {

	zaplog.Logger.Debugf("Connection close %s %v", c.RemoteAddr().String(), err)

	cc := c.Context().(*ctx.Context)
	service.RemRoomClients(cc.Broker, cc.RoomId, cc.ClientName)
	service.RemUserClient(cc.UserId, cc.ClientName)
	ctx.RemContext(c.RemoteAddr().String())
	cc = nil

	return
}

func (s *server) OnTraffic(c gnet.Conn) (action gnet.Action) {

	cc := c.Context().(*ctx.Context)
	var packets [][]byte
	for {
		data, err := s.serializer.Decode(c)
		if err == serializer.ErrIncompletePacket {
			break
		}
		if err != nil {
			logging.Errorf("invalid packet: %v", err)
			return gnet.Close
		}
		packet, _ := s.serializer.DecodePacket(data)

		ret := process(cc, packet)

		if ret == nil {
			continue
		}

		retBuf, _ := s.serializer.EncodePacket(ret)
		buf, _ := s.serializer.Encode(retBuf)

		packets = append(packets, buf)
	}
	if n := len(packets); n > 1 {
		_, _ = c.Writev(packets)
	} else if n == 1 {
		_, _ = c.Write(packets[0])
	}
	return
}

func process(c *ctx.Context, packet *protocol.Packet) *protocol.Packet {

	var ret *protocol.Packet = protocol.NewResponseError(packet, err.Default)
	var e error = nil
	switch packet.Header.Command {
	case protocol.CommandDefault:
		ret, e = handler.SingleDefaultHandler().Handle(c, packet)
		break
	case protocol.CommandSignal:
		ret, e = handler.SingleSignalHandler().Handle(c, packet)
		break
	case protocol.CommandContent:
		ret, e = handler.SingleContentHandler().Handle(c, packet)
		break
	case protocol.CommandCustom:
		ret, e = handler.CustomContentHandler().Handle(c, packet)
		break
	}

	if ret == nil {
		if e != nil {
			ret = protocol.NewResponseError(packet, err.Default)
		}
	}

	return ret
}

func getBrokerAddress() string {

	addr := ":" + config.OP.Port

	var brokerAddress string
	if util.IsNotEmpty(config.OP.Ip) {
		brokerAddress = config.OP.Ip + addr
	} else {
		brokerAddress = util.GetBrokerIp() + addr
	}
	return brokerAddress
}

func Start() {

	zaplog.InitLogger()

	zaplog.Logger.Infof("Start ...")

	brokerAddress := getBrokerAddress()

	ctx.BrokerAddress = brokerAddress

	service.SetBrokerInstance(brokerAddress)
	zaplog.Logger.Infof("Start broker address %s", brokerAddress)

	addr := "tcp://:" + config.OP.Port
	server := &server{
		serializer: serializer.NewJsonSerializer(),
	}
	op1 := gnet.WithMulticore(true)
	op2 := gnet.WithNumEventLoop(50)

	go startTask(brokerAddress)
	zaplog.Logger.Infof("Start scheduler start... %s", brokerAddress)

	if err := gnet.Run(server, addr, op1, op2); err != nil {
		util.Panic(err)
	}

	zaplog.Logger.Infof("Start Listen %s", brokerAddress)

}

func startTask(broker string) {
	c := cron.New()

	c.AddFunc("@every 1s", func() {
		service.SetBrokerInstance(broker)
		zaplog.Logger.Infof("Probe context size %d", ctx.ConnCount())

	})

	c.AddFunc("@every 1m", func() {
		service.RangeRoomClients("100", func(i interface{}) bool {
			fmt.Println("-------------- ", i)
			return true
		})
	})

	c.Start()
}
