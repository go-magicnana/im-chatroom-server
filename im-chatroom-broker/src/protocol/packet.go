package protocol

import (
	"encoding/json"
	"im-chatroom-broker/context"
	err "im-chatroom-broker/error"
)

const (
	MetaVersionBytes = 1
	MetaLengthBytes  = 4

	TargetAll = 1
	TargetOne = 2

	FlowUp   = 1
	FlowDown = 2

	CommandDefault = 1
	CommandSignal  = 2
	CommandNotice  = 3
	CommandContent = 4
	CommandGift    = 5
	CommandGoods   = 6

	TypeSignalPing       = 2101
	TypeSignalLogin      = 2102
	TypeSignalJoinRoom   = 2104
	TypeSignalChangeRoom = 2104
	TypeSignalLeaveRoom  = 2105

	TypeNoticeJoinRoom    = 3101
	TypeNoticeLeaveRoom   = 3102
	TypeNoticeBlockUser   = 3103
	TypeNoticeUnblockUser = 3104
	TypeNoticeCloseRoom   = 3105

	TypeContentText  = 4105
	TypeContentEmoji = 4105
	TypeContentReply = 4105

	TypeGiftNone = 5101

	TypeGoodsNone = 6101
)

type Packet struct {
	Header MessageHeader `json:"header"`
	Body   any           `json:"body"`
}

type MessageHeader struct {
	MessageId string `json:"messageId"`
	Command   uint16 `json:"command"`
	Target    uint32 `json:"target"`
	From      User   `json:"from"`
	To        User   `json:"to"`
	Flow      uint8  `json:"flow"`
	Type      uint32 `json:"type"`
	Code      uint32 `json:"code"`
	Message   string `json:"message"`
}

func NewResponseOK(in *Packet, body any) *Packet {

	header := MessageHeader{
		MessageId: in.Header.MessageId,
		Command:   in.Header.Command,
		Target:    in.Header.Target,
		From:      in.Header.From,
		To:        in.Header.To,
		Type:      in.Header.Type,
		Flow:      FlowDown,
		Code:      err.OK.Code,
		Message:   err.OK.Message,
	}

	return &Packet{
		Header: header,
		Body:   body,
	}
}

func NewResponseError(in *Packet, error err.Error) *Packet {
	header := MessageHeader{
		MessageId: in.Header.MessageId,
		Command:   in.Header.Command,
		Target:    in.Header.Target,
		From:      in.Header.From,
		To:        in.Header.To,
		Type:      in.Header.Type,
		Flow:      FlowDown,
		Code:      error.Code,
		Message:   error.Message,
	}

	return &Packet{
		Header: header,
	}
}

type User struct {
	Token  string `json:"token"`
	UserId string `json:"userId"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Role   string `json:"role"`
	RoomId string `json:"roomId"`
	Broker string `json:"broker"`
}

type MessageBodySignalPing struct {
	Current uint64 `json:"current"`
	UserId  string `json:"userId"`
}

func JsonSignalPing(any any, c *context.Context) *MessageBodySignalPing {
	bs, _ := json.Marshal(any)
	ret := MessageBodySignalPing{}
	json.Unmarshal(bs, &ret)
	return &ret
}

type MessageBodySignalLogin struct {
	Token string `json:"token"`
}

func JsonSignalLogin(any any, c *context.Context) *MessageBodySignalLogin {
	bs, _ := json.Marshal(any)
	ret := MessageBodySignalLogin{}
	json.Unmarshal(bs, &ret)
	return &ret
}

type MessageBodySignalDisconnect struct {
	UserId string `json:"userId"`
}

func JsonSignalDisconnect(any any, c *context.Context) *MessageBodySignalDisconnect {
	bs, _ := json.Marshal(any)
	ret := MessageBodySignalDisconnect{}
	json.Unmarshal(bs, &ret)
	return &ret
}

type MessageBodySignalJoinRoom struct {
	UserId string `json:"userId"`
	RoomId string `json:"roomId"`
}

func JsonSignalJoinRoom(any any, c *context.Context) *MessageBodySignalJoinRoom {
	bs, _ := json.Marshal(any)
	ret := MessageBodySignalJoinRoom{}
	json.Unmarshal(bs, &ret)
	return &ret
}

func JsonSignalLeaveRoom(any any, c *context.Context) *MessageBodySignalLeaveRoom {
	bs, _ := json.Marshal(any)
	ret := MessageBodySignalLeaveRoom{}
	json.Unmarshal(bs, &ret)
	return &ret
}

type MessageBodySignalLeaveRoom struct {
	UserId string `json:"userId"`
	RoomId string `json:"roomId"`
}

type MessageBodyNoticeJoinRoom struct {
	UserId string `json:"userId"`
	RoomId string `json:"roomId"`
}

type MessageBodyNoticeLeaveRoom struct {
	UserId string `json:"userId"`
	RoomId string `json:"roomId"`
}

type MessageBodyNoticeBlockUser struct {
	UserId string `json:"userId"`
	RoomId string `json:"roomId"`
}

type MessageBodyNoticeUnblockUser struct {
	UserId string `json:"userId"`
	RoomId string `json:"roomId"`
}

type MessageBodyNoticeCloseRoom struct {
	RoomId string `json:"roomId"`
}

type MessageBodyNoticeBlockRoom struct {
	RoomId string `json:"roomId"`
}

type MessageBodyNoticeUnblockRoom struct {
	RoomId string `json:"roomId"`
}

type MessageBodyContentText struct {
	Content string `json:"content"`
}

type MessageBodyContentEmoji struct {
	Content string `json:"content"`
}

type MessageBodyContentReply struct {
	SendMessageId  string `json:"sendMessageId"`
	SendUserId     string `json:"sendUserId"`
	SendUserName   string `json:"sendUserName"`
	SendUserAvatar string `json:"sendUserAvatar"`
	SendContent    string `json:"sendContent"`
	Content        string `json:"content"`
}
