package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-uuid"
	"github.com/labstack/echo"
	"github.com/ziflex/lecho/v3"
	"golang.org/x/net/context"
	"im-chatroom-gateway/src/mq"
	"im-chatroom-gateway/src/protocol"
	"im-chatroom-gateway/src/service"
	"net/http"
	"os"
	"strconv"
)

func MessagePush(ct echo.Context) error {
	e := echo.New()
	e.Logger = lecho.New(
		os.Stdout,
		lecho.WithFields(map[string]interface{}{"name": "lecho factory"}),
		lecho.WithTimestamp(),
		lecho.WithCaller(),
		lecho.WithPrefix("controllers.MessagePush"),
	)

	//获取post请求的表单参数，
	// 类型 是何种通知
	fromUserId := ct.FormValue("fromUserId")
	fromUserName := ct.FormValue("fromUserName")
	fromUserAvatar := ct.FormValue("fromUserAvatar")
	// 类型 是何种通知
	messageTarget := ct.FormValue("messageTarget")
	messageTargetInt64, _ := strconv.ParseInt(messageTarget, 0, 64)

	messageType := ct.FormValue("messageType")
	messageTypeInt64, _ := strconv.ParseInt(messageType, 10, 64)

	// 个人 userid，所有人 roomid
	userId := ct.FormValue("userId")
	// 信息
	roomId := ct.FormValue("roomId")

	messageId, _ := uuid.GenerateUUID()

	//TypeNoticeBlockUser   = 3103
	//TypeNoticeUnblockUser = 3104
	//TypeNoticeCloseRoom   = 3105
	//TypeNoticeBlockRoom   = 3106
	//TypeNoticeUnblockRoom = 3107

	header := protocol.MessageHeader{
		MessageId: messageId,
		Command:   protocol.CommandNotice,
		Target:    uint32(messageTargetInt64),
		From: protocol.UserInfo{
			UserId: fromUserId,
			Name:   fromUserName,
			Avatar: fromUserAvatar,
		},
		To:   userId,
		Flow: protocol.FlowUp,
		Type: uint32(messageTypeInt64),
	}

	var body any

	switch int(messageTypeInt64) {
	case protocol.TypeNoticeBlockUser:
		body = protocol.MessageBodyNoticeBlockUser{UserId: userId, RoomId: roomId}

	case protocol.TypeNoticeUnblockUser:
		body = protocol.MessageBodyNoticeUnblockUser{UserId: userId, RoomId: roomId}

	case protocol.TypeNoticeCloseRoom:
		body = protocol.MessageBodyNoticeCloseRoom{RoomId: roomId}

	case protocol.TypeNoticeBlockRoom:
		body = protocol.MessageBodyNoticeBlockRoom{RoomId: roomId}

	case protocol.TypeNoticeUnblockRoom:
		body = protocol.MessageBodyNoticeUnblockRoom{RoomId: roomId}

	}

	packet := protocol.Packet{
		Header: header, Body: body,
	}

	fmt.Println("send message packet :" + packet.ToString())

	e.Logger.Info("send notice message ", packet)

	result := deliver(context.Background(), &packet)
	if result != nil {
		e.Logger.Info("send notice message error:", result)
		return ct.JSON(http.StatusOK, gin.H{"code": 1001, "message": "Server Error"})
	}

	return ct.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}

func deliver(ctx context.Context, packet *protocol.Packet) error {

	packet.Header.Flow = protocol.FlowDeliver

	if packet.Header.Target == protocol.TargetRoom {
		mq.SendSync2Room(packet)
	} else {

		ret := service.GetUserClients(ctx, packet.Header.To)

		for _, v := range ret {

			msg := &protocol.PacketMessage{
				UserKey: v,
				Packet:  *packet,
			}

			broker, _ := service.GetUserDeviceBroker(ctx, v)

			mq.SendSync2One(broker, msg)
			//fmt.Println(msg,broker)
		}

	}

	return nil
}
