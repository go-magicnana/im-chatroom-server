package controllers

import (
	"github.com/hashicorp/go-uuid"
	"github.com/labstack/echo"
	"golang.org/x/net/context"
	"im-chatroom-gateway/apierror"
	"im-chatroom-gateway/redis"
	"im-chatroom-gateway/zaplog"
	"math/rand"
	"net/http"
	"time"
)

const (
	RoomInfo string = "imchatroom:roominfo:"
)

func CreateChatroom(c echo.Context) error {

	zaplog.Logger.Debugf("%s %v", c.Request().RequestURI, nil)

	id, _ := uuid.GenerateUUID()

	return write(c, http.StatusOK, NewApiResultOK(id))

}

func randCreator(l int) string {
	str := "0123456789abcdefghigklmnopqrstuvwxyz"
	strList := []byte(str)

	result := []byte{}
	i := 0

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i < l {
		new := strList[r.Intn(len(strList))]
		result = append(result, new)
		i = i + 1
	}
	return string(result)
}

func GetRoomMembers(ct echo.Context) error {

	//获取post请求的表单参数
	roomId := ct.FormValue("roomId")
	if roomId == "" {
		e.Logger.Info("roomId is illegal")
		return ct.JSON(http.StatusOK, NewApiResultError(apierror.InvalidParameter))
	}

	snum, err := redis.Rdb.SCard(context.Background(), RoomInfo+roomId).Result()
	if err == redis.Nil {
		snum = 0
	}

	return ct.JSON(http.StatusOK, NewApiResultOK(snum))
}
