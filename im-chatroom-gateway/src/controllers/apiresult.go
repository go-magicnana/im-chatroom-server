package controllers

import "im-chatroom-gateway/apierror"

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
