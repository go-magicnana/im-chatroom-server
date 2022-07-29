package controllers

import (
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/labstack/echo"
	"github.com/ziflex/lecho/v3"
	"golang.org/x/net/context"
	"im-chatroom-gateway/apierror"
	"im-chatroom-gateway/domains"
	"im-chatroom-gateway/protocol"
	"im-chatroom-gateway/service"
	"im-chatroom-gateway/zaplog"
	"net/http"
	"os"
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

func MessagePush(c echo.Context) error {

	/**
	TypeNoticeBlockUser   = 3103
	TypeNoticeUnblockUser = 3104
	TypeNoticeCloseRoom   = 3105
	TypeNoticeBlockRoom   = 3106
	TypeNoticeUnblockRoom = 3107
	*/

	a, _ := c.FormParams()
	zaplog.Logger.Debugf("%s %v", c.Request().RequestURI, a)

	if len(a) == 0 {
		return write(c, http.StatusOK, NewApiResultError(apierror.InvalidParameter))
	}

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
			Target:    protocol.TargetOne,
			To:        u.UserId,
			Flow:      protocol.FlowDeliver,
			Type:      protocol.TypeNoticeBlockUser,
		},
		Body: protocol.MessageBodyNoticeBlockUser{
			User:   *userInfo,
			RoomId: u.RoomId,
		},
	}
	service.SetRoomMemberBlocked(context.Background(), u.RoomId, u.UserId)
	SendSync2User(context.Background(), &packet)

	return write(c, http.StatusOK, NewApiResultOK(nil))

}

