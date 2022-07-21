package controllers

import (
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/labstack/echo"
	"github.com/ziflex/lecho/v3"
	"golang.org/x/net/context"
	"im-chatroom-gateway/apierror"
	"im-chatroom-gateway/mq"
	"im-chatroom-gateway/protocol"
	"im-chatroom-gateway/service"
	"net/http"
	"os"
	"strconv"
)

var e echo.Echo

func init() {
	fmt.Println("messagecontroller init")
	e := echo.New()
	e.Logger = lecho.New(
		os.Stdout,
		lecho.WithFields(map[string]interface{}{"name": "lecho factory"}),
		lecho.WithTimestamp(),
		lecho.WithCaller(),
		lecho.WithPrefix("controllers.MessageController"),
	)
}

func MessagePush(ct echo.Context) error {

	/**
	TypeNoticeBlockUser   = 3103
	TypeNoticeUnblockUser = 3104
	TypeNoticeCloseRoom   = 3105
	TypeNoticeBlockRoom   = 3106
	TypeNoticeUnblockRoom = 3107
	*/

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
		Flow: protocol.FlowDeliver,
		Type: uint32(messageTypeInt64),
		Code: 200,
	}

	var body any
	var userinfo *protocol.UserInfo
	switch int(messageTypeInt64) {
	case protocol.TypeNoticeBlockUser:
		userinfo, _ = service.GetUserInfo(context.Background(), userId)
		body = protocol.MessageBodyNoticeBlockUser{User: *userinfo, RoomId: roomId}

	case protocol.TypeNoticeUnblockUser:
		userinfo, _ = service.GetUserInfo(context.Background(), userId)
		body = protocol.MessageBodyNoticeUnblockUser{User: *userinfo, RoomId: roomId}

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

	e.Logger.Info("send notice message ", packet)

	result := deliver(context.Background(), &packet, roomId)
	if result != nil {
		e.Logger.Info("send notice message error:", result)
		return ct.JSON(http.StatusOK, NewApiResultError(apierror.Default))
	}

	return ct.JSON(http.StatusOK, NewApiResultOK(nil))
}

func deliver(ctx context.Context, packet *protocol.Packet, roomId string) error {

	packet.Header.Flow = protocol.FlowDeliver

	if packet.Header.Target == protocol.TargetRoom {
		packet.Header.To = roomId
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
