package error

var OK = Error{Code: 200, Message: "OK"}
var Default = Error{Code: 1001, Message: "Server Error"}
var CommandNotAllow = Error{1002, "Command Not Allow"}
var TypeNotAllow = Error{1003, "Type Not Allow"}
var Unauthorized = Error{1004, "Unauthorized"}

type Error struct {
	Code    uint32
	Message string
}
