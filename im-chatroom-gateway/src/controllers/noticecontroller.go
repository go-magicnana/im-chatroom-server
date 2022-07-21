package controllers

import (
	"github.com/hashicorp/go-uuid"
	"github.com/labstack/echo"
	"golang.org/x/net/context"
	"im-chatroom-gateway/apierror"
	"im-chatroom-gateway/domains"
	"im-chatroom-gateway/mq"
	"im-chatroom-gateway/protocol"
	"im-chatroom-gateway/service"
	"im-chatroom-gateway/util"
	"im-chatroom-gateway/zaplog"
	"net/http"
)

func NoticeBlockUser(c echo.Context) error {

	/**
	TypeNoticeBlockUser   = 3103
	TypeNoticeUnblockUser = 3104
	TypeNoticeCloseRoom   = 3105
	TypeNoticeBlockRoom   = 3106
	TypeNoticeUnblockRoom = 3107
	*/

	a, _ := c.FormParams()
	zaplog.Logger.Debugf("%s %v", c.Request().RequestURI, a)

	u := new(domains.BlockUser)

	if err := c.Bind(u); err != nil {
		return write(c, http.StatusOK, NewApiResultError(err))
	}

	if err := c.Validate(u); err != nil {
		return write(c, http.StatusOK, NewApiResultError(err))
	}

	userInfo, err := service.GetUserInfo(context.Background(), u.UserId)
	if err != nil {
		return write(c, http.StatusOK, NewApiResultError(err))
	}

	msgId, _ := uuid.GenerateUUID()

	packet := protocol.Packet{
		Header: protocol.MessageHeader{
			MessageId: msgId,
			Command:   protocol.CommandNotice,
			Target:    protocol.TargetRoom,
			To:        u.RoomId,
			Flow:      protocol.FlowDeliver,
			Type:      protocol.TypeNoticeBlockUser,
		},
		Body: protocol.MessageBodyNoticeBlockUser{
			User:   *userInfo,
			RoomId: u.RoomId,
		},
	}

	mq.SendSync2Room(&packet)

	return write(c, http.StatusOK, NewApiResultOK(nil))

}

func NoticeUnblockUser(c echo.Context) error {

	/**
	TypeNoticeBlockUser   = 3103
	TypeNoticeUnblockUser = 3104
	TypeNoticeCloseRoom   = 3105
	TypeNoticeBlockRoom   = 3106
	TypeNoticeUnblockRoom = 3107
	*/

	a, _ := c.FormParams()
	zaplog.Logger.Debugf("%s %v", c.Request().RequestURI, a)

	u := new(domains.BlockUser)

	if err := c.Bind(u); err != nil {
		return write(c, http.StatusOK, NewApiResultError(err))
	}

	if err := c.Validate(u); err != nil {
		return write(c, http.StatusOK, NewApiResultError(err))
	}

	userInfo, err := service.GetUserInfo(context.Background(), u.UserId)
	if err != nil {
		return write(c, http.StatusOK, NewApiResultError(err))
	}

	msgId, _ := uuid.GenerateUUID()

	packet := protocol.Packet{
		Header: protocol.MessageHeader{
			MessageId: msgId,
			Command:   protocol.CommandNotice,
			Target:    protocol.TargetRoom,
			To:        u.RoomId,
			Flow:      protocol.FlowDeliver,
			Type:      protocol.TypeNoticeUnblockUser,
		},
		Body: protocol.MessageBodyNoticeUnblockUser{
			User:   *userInfo,
			RoomId: u.RoomId,
		},
	}

	mq.SendSync2Room(&packet)

	return write(c, http.StatusOK, NewApiResultOK(nil))

}

func NoticeCloseRoom(c echo.Context) error {

	/**
	TypeNoticeBlockUser   = 3103
	TypeNoticeUnblockUser = 3104
	TypeNoticeCloseRoom   = 3105
	TypeNoticeBlockRoom   = 3106
	TypeNoticeUnblockRoom = 3107
	*/

	a, _ := c.FormParams()
	zaplog.Logger.Debugf("%s %v", c.Request().RequestURI, a)

	roomId := c.FormValue("roomId")
	if util.IsEmpty(roomId) {
		return write(c, http.StatusOK, NewApiResultError(apierror.InvalidParameter))
	}

	msgId, _ := uuid.GenerateUUID()

	packet := protocol.Packet{
		Header: protocol.MessageHeader{
			MessageId: msgId,
			Command:   protocol.CommandNotice,
			Target:    protocol.TargetRoom,
			To:        roomId,
			Flow:      protocol.FlowDeliver,
			Type:      protocol.TypeNoticeCloseRoom,
		},
		Body: protocol.MessageBodyNoticeCloseRoom{
			RoomId: roomId,
		},
	}

	mq.SendSync2Room(&packet)

	return write(c, http.StatusOK, NewApiResultOK(nil))

}

func NoticeBlockRoom(c echo.Context) error {

	/**
	TypeNoticeBlockUser   = 3103
	TypeNoticeUnblockUser = 3104
	TypeNoticeCloseRoom   = 3105
	TypeNoticeBlockRoom   = 3106
	TypeNoticeUnblockRoom = 3107
	*/

	a, _ := c.FormParams()
	zaplog.Logger.Debugf("%s %v", c.Request().RequestURI, a)

	roomId := c.FormValue("roomId")
	if util.IsEmpty(roomId) {
		return write(c, http.StatusOK, NewApiResultError(apierror.InvalidParameter))
	}

	msgId, _ := uuid.GenerateUUID()

	packet := protocol.Packet{
		Header: protocol.MessageHeader{
			MessageId: msgId,
			Command:   protocol.CommandNotice,
			Target:    protocol.TargetRoom,
			To:        roomId,
			Flow:      protocol.FlowDeliver,
			Type:      protocol.TypeNoticeBlockRoom,
		},
		Body: protocol.MessageBodyNoticeBlockRoom{
			RoomId: roomId,
		},
	}

	mq.SendSync2Room(&packet)

	return write(c, http.StatusOK, NewApiResultOK(nil))

}

func NoticeUnblockRoom(c echo.Context) error {

	/**
	TypeNoticeBlockUser   = 3103
	TypeNoticeUnblockUser = 3104
	TypeNoticeCloseRoom   = 3105
	TypeNoticeBlockRoom   = 3106
	TypeNoticeUnblockRoom = 3107
	*/

	a, _ := c.FormParams()
	zaplog.Logger.Debugf("%s %v", c.Request().RequestURI, a)

	roomId := c.FormValue("roomId")
	if util.IsEmpty(roomId) {
		return write(c, http.StatusOK, NewApiResultError(apierror.InvalidParameter))
	}

	msgId, _ := uuid.GenerateUUID()

	packet := protocol.Packet{
		Header: protocol.MessageHeader{
			MessageId: msgId,
			Command:   protocol.CommandNotice,
			Target:    protocol.TargetRoom,
			To:        roomId,
			Flow:      protocol.FlowDeliver,
			Type:      protocol.TypeNoticeUnblockRoom,
		},
		Body: protocol.MessageBodyNoticeUnblockRoom{
			RoomId: roomId,
		},
	}

	mq.SendSync2Room(&packet)

	return write(c, http.StatusOK, NewApiResultOK(nil))

}
