package controllers

import (
	"github.com/hashicorp/go-uuid"
	"github.com/labstack/echo"
	"im-chatroom-gateway/zaplog"
	"math/rand"
	"net/http"
	"time"
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
