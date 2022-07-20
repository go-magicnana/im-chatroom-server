package apierror

import "fmt"

var OK = ApiError{Code: 200, Message: "OK"}
var Default = ApiError{Code: 1001, Message: "Server Error"}
var InvalidParameter = ApiError{1002, "Invalid Parameter"}
var CouldNotBeSeries = ApiError{1003, "Series Error %s"}
var StorageResponseNil = ApiError{1004, "Storage response nil"}
var StorageResponseError = ApiError{1004, "Storage response error %s"}

type ApiError struct {
	Code    uint32
	Message string
}

func (e ApiError) WrapperAndFormat(error error) error{
	e.Message = fmt.Sprintf(e.Message,error.Error())
	return e
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
