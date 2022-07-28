package client

import (
	"bytes"
	"encoding/binary"
	"github.com/hashicorp/go-uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"im-chatroom-gateway/protocol"
	"im-chatroom-gateway/serializer"
	"im-chatroom-gateway/service"
	"im-chatroom-gateway/zaplog"
	"io"
	"net"
	"sync"
	"time"
)

var brokers sync.Map

func Heartbeat() {

	c := context.Background()

	queryRedisAndStartHeartBeat(c)

}

func queryRedisAndStartHeartBeat(ctx context.Context) {
	for {
		brokers := service.GetBrokerInstance(ctx)

		if brokers != nil && len(brokers) > 0 {

			for _, broker := range brokers {

				root, beatCancel := context.WithCancel(ctx)

				go doHeartbeat(root, beatCancel, broker)
			}
		}
	}
}

func doHeartbeat(c context.Context, cancel context.CancelFunc, broker string) {

	if broker == "" {
		return
	}

	v, b := brokers.Load(broker)

	if v == "OK" || b {
		return
	}

	zaplog.Logger.Infof("Heartbeat %s Start", broker)
	brokers.Store(broker, "OK")
	connect(c, cancel, broker)
	zaplog.Logger.Infof("Heartbeat %s Quit", broker)

}

func connect(c context.Context, cancel context.CancelFunc, broker string) {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", broker)
	if err != nil {
		zaplog.Logger.Debugf("Heartbeat %s ResloveError", broker)
		clearBroker(c, broker)
		cancel()
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		zaplog.Logger.Debugf("Heartbeat %s ConnectError", broker)
		clearBroker(c, broker)
		cancel()
		return
	}

	zaplog.Logger.Debugf("Heartbeat %s StartBeat", broker)

	ch := make(chan string, 1)

	if err := sendHeartBeat(conn); err != nil {
		close(c, cancel, broker, conn)
	}

	doReadAndHearBeat(c, cancel, broker, conn, ch)

}

func doReadAndHearBeat(c context.Context, cancel context.CancelFunc, broker string, conn net.Conn, ch chan string) {

	defer close(c, cancel, broker, conn)

	go doRead(c, ch, conn)
	go doHeartBeat(c, conn, ch)
}

func doHeartBeat(c context.Context, conn net.Conn, ch chan string) {
	for {
		select {
		case <-c.Done():
			return
		case body := <-ch:
			if "QUIT" == body {
				return
			} else {
				time.Sleep(time.Second * 5)
				if err := sendHeartBeat(conn); err != nil {
					return
				} else {
					continue
				}
			}
		default:
			continue
		}
	}
}

func doRead(c context.Context, ch chan string, conn net.Conn) {

	for {
		select {
		case <-c.Done():
			return
		default:

			serializer := serializer.SingleJsonSerializer()

			conn.SetReadDeadline(time.Now().Add(time.Second * 10))

			meta := make([]byte, 5)
			ml, me := conn.Read(meta)

			switch me.(type) {
			case *net.OpError:
				zaplog.Logger.Errorf("Heartbeat %s ReadTimeOut", conn.RemoteAddr().String())
				ch <- "QUIT"
				return
			}

			if me == io.EOF {

				zaplog.Logger.Errorf("Heartbeat %s ReadClose", conn.RemoteAddr().String())
				ch <- "QUIT"
				return
			}

			if me != nil {
				zaplog.Logger.Errorf("Heartbeat %s ReadError", conn.RemoteAddr().String())
				ch <- "QUIT"
				return
			}

			if ml != 5 {
				zaplog.Logger.Errorf("Heartbeat %s MetaError", conn.RemoteAddr().String())
				ch <- "QUIT"
				return
			}

			version := meta[0]

			if version != serializer.Version() {
				zaplog.Logger.Errorf("Heartbeat %s VersionError", conn.RemoteAddr().String())
				ch <- "QUIT"
				return
			}

			length := binary.BigEndian.Uint32(meta[1:])
			body := make([]byte, length)
			conn.Read(body)

			p, _ := serializer.DecodePacket(body)

			zaplog.Logger.Debugf("Heartbeat %s ReadOK %s %d %d %s", conn.RemoteAddr().String(), p.Header.MessageId, p.Header.Command, p.Header.Type, p.Body)

			ch <- p.Body.(string)

		}
	}
}

func close(c context.Context, cancel context.CancelFunc, broker string, conn net.Conn) {
	cancel()
	clearBroker(c, broker)
	conn.Close()
	zaplog.Logger.Infof("Heartbeat %s cancel clear close", broker)

}

func clearBroker(ctx context.Context, broker string) {
	zaplog.Logger.Infof("Heartbeat %s ClearBroker", broker)

	service.DelBrokerInstance(ctx, broker)
	service.DelBrokerCapacityAll(ctx, broker)
}

func sendHeartBeat(conn net.Conn) error {
	return write(conn, heartBeatPacket())
}

func heartBeatPacket() *protocol.Packet {
	uuid, _ := uuid.GenerateUUID()

	header := protocol.MessageHeader{
		MessageId: uuid,
		Command:   protocol.CommandDefault,
		Flow:      protocol.FlowUp,
		Type:      protocol.TypeDefaultHeartBeat,
	}

	body := protocol.MessageBodyDefaultHeartBeat{
		Password: protocol.TypeDefaultHeartBeatPassword,
	}

	return &protocol.Packet{
		Header: header, Body: body,
	}
}

func write(conn net.Conn, p *protocol.Packet) error {

	j := serializer.SingleJsonSerializer()

	bs, e := j.EncodePacket(p)
	if bs == nil {
		return errors.New("empty packet")
	}

	if e != nil {
		return e
	}

	buffer := new(bytes.Buffer)

	binary.Write(buffer, binary.BigEndian, j.Version())

	length := uint32(len(bs))
	binary.Write(buffer, binary.BigEndian, length)

	buffer.Write(bs)
	_, err := conn.Write(buffer.Bytes())

	zaplog.Logger.Debugf("Heartbeat %s WriteOK %s %d %d %s", conn.RemoteAddr().String(), p.Header.MessageId, p.Header.Command, p.Header.Type, p.Body)

	if err != nil {
		return errors.New("write response error +" + err.Error())
	} else {
		return nil
	}

}

func read(conn net.Conn) *protocol.Packet {

	serializer := serializer.SingleJsonSerializer()

	conn.SetReadDeadline(time.Now().Add(time.Second * 10))

	meta := make([]byte, 5)
	ml, me := conn.Read(meta)

	switch me.(type) {
	case *net.OpError:
		zaplog.Logger.Errorf("Heartbeat %s ReadTimeOut", conn.RemoteAddr().String())
		return nil
	}

	if me == io.EOF {

		zaplog.Logger.Errorf("Heartbeat %s ReadClose", conn.RemoteAddr().String())
		return nil
	}

	if me != nil {
		zaplog.Logger.Errorf("Heartbeat %s ReadError", conn.RemoteAddr().String())
		return nil
	}

	if ml != 5 {
		zaplog.Logger.Errorf("Heartbeat %s MetaError", conn.RemoteAddr().String())
		return nil
	}

	version := meta[0]

	if version != serializer.Version() {
		zaplog.Logger.Errorf("Heartbeat %s VersionError", conn.RemoteAddr().String())
		return nil
	}

	length := binary.BigEndian.Uint32(meta[1:])
	body := make([]byte, length)
	conn.Read(body)

	p, _ := serializer.DecodePacket(body)

	zaplog.Logger.Debugf("Heartbeat %s ReadOK %s %d %d %s", conn.RemoteAddr().String(), p.Header.MessageId, p.Header.Command, p.Header.Type, p.Body)

	return p
}
