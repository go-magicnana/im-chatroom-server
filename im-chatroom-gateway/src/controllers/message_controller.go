package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo"
	"net/http"
)

func MessagePush(ct echo.Context) error {

	return ct.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}
