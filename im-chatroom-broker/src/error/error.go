package error

import "fmt"

var OK = Error{Code: 200, Message: "OK"}
var Default = Error{Code: 1001, Message: "Server Error"}


var InvalidRequest = Error{1001,"Invalid Request Parameter [%s]"}
var NoSuchDataExist = Error{1002,"No Such Data Exist"}

var CommandNotAllow = Error{1002, "Command Not Allow"}
var TypeNotAllow = Error{1003, "Type Not Allow"}
var Unauthorized = Error{1004, "Unauthorized"}
var AlreadyLogin = Error{1005, "Already Login"}

type Error struct {
	Code    uint32
	Message string
}

func (e Error) Format(msg string) Error{
	e.Message = fmt.Sprintf(e.Message,msg)
	return e
}
