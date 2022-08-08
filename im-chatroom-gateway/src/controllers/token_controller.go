package controllers

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"im-chatroom-gateway/apierror"
	"im-chatroom-gateway/domains"
	"im-chatroom-gateway/redis"
	"im-chatroom-gateway/zaplog"
	"math/rand"
	"net/http"
	"time"
)

const userTokenKey string = "imchatroom:user.auth:"

func GetToken(c echo.Context) error {

	a, _ := c.FormParams()
	zaplog.Logger.Debugf("%s %v", c.Request().RequestURI, a)

	if len(a) == 0 {
		return write(c, http.StatusOK, NewApiResultError(apierror.InvalidParameter))
	}

	u := new(domains.UserInfo)

	if err := c.Bind(u); err != nil {
		return write(c, http.StatusOK, NewApiResultError(err))
	}

	if err := c.Validate(u); err != nil {
		return write(c, http.StatusOK, NewApiResultError(err))
	}

	u.Token = buildToken(u.UserId)

	//u.Token = randCreator22(4)

	if err := SetUserAuth(*u); err != nil {
		return write(c, http.StatusOK, NewApiResultError(err))
	}

	return write(c, http.StatusOK, NewApiResultOK(u.Token))

}

func write(c echo.Context, code int, ret ApiResult) error {
	zaplog.Logger.Infof("%s %v", c.Request().RequestURI, ret)
	return c.JSON(code, ret)
}

func buildToken(userId string) string {
	timeUnix := time.Now().Unix()
	data := []byte("userToken:" + userId + fmt.Sprintf("%d", timeUnix))
	// 将[]byte转成16进制
	userToken := fmt.Sprintf("%x", md5.Sum(data))
	return userToken
}

func SetUserAuth(u domains.UserInfo) error {

	data, err := json.Marshal(u)

	if err != nil {
		return apierror.CouldNotBeSeries.Format(err.Error())
	}

	result := redis.Rdb.Set(context.Background(), userTokenKey+u.Token, data, -1)
	if result == nil {
		return apierror.StorageResponseNil
	}

	_, e := result.Result()
	if e != nil {
		return apierror.StorageResponseError.Format(e.Error())
	}

	return nil

}

func randCreator22(l int) string {
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
