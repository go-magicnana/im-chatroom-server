package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
	"im-chatroom-gateway/apierror"
)

type ApiResult struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func NewApiResultError(e error) ApiResult {

	switch e.(type) {
	case apierror.ApiError:
		apiError := e.(apierror.ApiError)
		return ApiResult{
			Code: apiError.Code,
			Msg:  apiError.Error(),
		}
	case validator.ValidationErrors:
		return ApiResult{
			Code: apierror.InvalidParameter.Code,
			Msg:  e.Error(),
		}

	case *echo.HTTPError:
		return ApiResult{
			Code:apierror.ParameterBindError.Code,
			Msg: e.Error(),
		}
	default:
		return ApiResult{
			Code: apierror.Default.Code,
			Msg:  e.Error(),
		}
	}

}

func NewApiResultOK(any any) ApiResult {
	return ApiResult{
		Code: apierror.OK.Code,
		Data: any,
	}
}
