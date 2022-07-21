package apierror

import (
	"fmt"
)

var OK = ApiError{Code: 200, Message: "OK"}
var Default = ApiError{Code: 1001, Message: "Server Error"}
var InvalidParameter = ApiError{1002, "Invalid Parameter"}
var ParameterBindError = ApiError{1003, "No parameters found"}
var CouldNotBeSeries = ApiError{1004, "Series Error %s"}
var StorageResponseNil = ApiError{1005, "Storage response nil"}
var StorageResponseError = ApiError{1006, "Storage response error %s"}
var StorageResponseEmpty = ApiError{1007, "Storage response empty"}

type ApiError struct {
	Code    uint32
	Message string
}

func (e ApiError) Format(msg string) error {
	e.Message = fmt.Sprintf(e.Message, msg)
	return e
}

func (e ApiError) Replace(msg string) error {
	e.Message = msg
	return e
}

func (e ApiError) Error() string {
	return e.Message
}


