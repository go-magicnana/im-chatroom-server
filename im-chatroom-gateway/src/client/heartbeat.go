package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"im-chatroom-gateway/protocol"
	"im-chatroom-gateway/serializer"
	"im-chatroom-gateway/service"
	"im-chatroom-gateway/zaplog"
	"io"
	"net"
	"os"
	"time"
)

func Heartbeat() {

	c := context.Background()

	brokers := service.GetBrokerInstance(c)

	if brokers != nil && len(brokers) > 0 {

		for _, broker := range brokers {

			go doHeartbeat(c, broker)
		}
	}

}

func doHeartbeat(c context.Context, broker string) {
	if service.Lock(c, broker) {
		defer service.Unlock(c, broker)
		connect(c,broker)
	}
}

func connect(c context.Context,broker string) {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", broker)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		clearBroker(c,broker)
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		clearBroker(c,broker)
		return
	}

	zaplog.Logger.Infof("Heartbeat %s start ...", broker)

	for {
		time.Sleep(time.Second)
		sendHeartBeat(broker, conn)
	}

}

func close(broker string, conn net.Conn) {
	clearBroker(context.Background(),broker)
	conn.Close()
}

func clearBroker(ctx context.Context,broker string){
	service.DelBrokerInstance(ctx,broker)
	service.DelBrokerCapacityAll(ctx,broker)
}

func sendHeartBeat(broker string, conn net.Conn) {

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

	packet := protocol.Packet{
		Header: header, Body: body,
	}

	write(conn, &packet)

	ret := read(conn)

	if ret.Body.(string) == "OK" {

	} else {
		close(broker, conn)
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

	zaplog.Logger.Debugf("WriteOK %s %s %d %d %s", conn.RemoteAddr().String(), p.Header.MessageId, p.Header.Command, p.Header.Type, p.Body)

	if err != nil {
		return errors.New("write response error +" + err.Error())
	} else {
		return nil
	}

}

func read(conn net.Conn) *protocol.Packet {

	serializer := serializer.SingleJsonSerializer()

	conn.SetReadDeadline(time.Now().Add(time.Second * 3))

	meta := make([]byte, 5)
	ml, me := conn.Read(meta)

	switch me.(type) {
	case *net.OpError:
		zaplog.Logger.Errorf("ReadTimeOut %s Close Client", conn.RemoteAddr().String())
		return nil
	}

	if me == io.EOF {

		zaplog.Logger.Errorf("ReadClose %s Close Client", conn.RemoteAddr().String())
		return nil
	}

	if me != nil {

		zaplog.Logger.Errorf("ReadError %s To Read Continue", conn.RemoteAddr().String())
		return nil
	}

	if ml != 5 {
		zaplog.Logger.Errorf("MetaError %s To Read Continue", conn.RemoteAddr().String())
		return nil
	}

	version := meta[0]

	if version != serializer.Version() {
		zaplog.Logger.Errorf("MetaOfVersionError %s To Read Continue", conn.RemoteAddr().String())
		return nil
	}

	length := binary.BigEndian.Uint32(meta[1:])
	body := make([]byte, length)
	conn.Read(body)

	packet, _ := serializer.DecodePacket(body)

	zaplog.Logger.Debugf("ReadOK %s Go Process %s %d %d %s", conn.RemoteAddr().String(), packet.Header.MessageId, packet.Header.Command, packet.Header.Type, packet.Body)

	return packet
}
