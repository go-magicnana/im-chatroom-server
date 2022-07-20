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
	"net/http"
	"time"
)

const userTokenKey string = "imchatroom:userauth:"

func GetToken(c echo.Context) error {

	u := new(domains.UserInfo)

	if err := c.Bind(u); err != nil {
		fmt.Println(err)
		return nil
	}

	u.Token = buildToken(u.UserId)

	_, e := SetUserAuth(*u)

	if e != nil {
		c.JSON(http.StatusOK, NewApiResultError(e))
	}

	return c.JSON(http.StatusOK, NewApiResultOK(u.Token))
}

func buildToken(userId string) string {
	timeUnix := time.Now().Unix()
	data := []byte("userToken:" + userId + fmt.Sprintf("%d", timeUnix))
	// 将[]byte转成16进制
	userToken := fmt.Sprintf("%x", md5.Sum(data))
	return userToken
}

func SetUserAuth(u domains.UserInfo) (string, error) {

	data, err := json.Marshal(u)

	if err != nil {
		return "", apierror.CouldNotBeSeries.WrapperAndFormat(err)
	}

	result := redis.Rdb.Set(context.Background(), userTokenKey+u.Token, data, time.Minute*30)
	if result == nil {
		return "", apierror.StorageResponseNil
	}

	ret, e := result.Result()
	if e != nil {
		return "", apierror.StorageResponseError.WrapperAndFormat(e)
	}

	return ret, nil

}
