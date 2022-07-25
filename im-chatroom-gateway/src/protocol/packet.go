package protocol

import (
	"encoding/json"
)

const (
	MetaVersionBytes = 1
	MetaLengthBytes  = 4

	TargetRoom = 1
	TargetOne  = 2

	FlowUp              = 1
	FlowDown            = 2
	FlowDeliver         = 3
	FlowDeliverResponse = 4

	CommandDefault = 1
	CommandSignal  = 2
	CommandNotice  = 3
	CommandContent = 4
	CommandGift    = 5
	CommandGoods   = 6

	TypeSignalPing       = 2101
	TypeSignalLogin      = 2102
	TypeSignalJoinRoom   = 2104
	TypeSignalLeaveRoom  = 2105
	TypeSignalChangeRoom = 2106

	TypeNoticeJoinRoom    = 3101
	TypeNoticeLeaveRoom   = 3102
	TypeNoticeBlockUser   = 3103
	TypeNoticeUnblockUser = 3104
	TypeNoticeCloseRoom   = 3105
	TypeNoticeBlockRoom   = 3106
	TypeNoticeUnblockRoom = 3107

	TypeContentText  = 4101
	TypeContentEmoji = 4102
	TypeContentAt    = 4103
	TypeContentReply = 4104

	TypeGiftNone = 5101

	TypeGoodsNone = 6101
)

type Packet struct {
	Header MessageHeader `json:"header"`
	Body   any           `json:"body"`
}

func (p Packet) ToString() string {
	bs, _ := json.Marshal(p)
	return string(bs)
}

type PacketMessage struct {
	ClientName string `json:"clientName"`
	Packet     Packet `json:"packet"`
}

type MessageHeader struct {
	MessageId string   `json:"messageId"`
	Command   uint16   `json:"command"`
	Target    uint32   `json:"target"`
	From      UserInfo `json:"from"`
	To        string   `json:"to"`
	Flow      uint8    `json:"flow"`
	Type      uint32   `json:"type"`
	Code      uint32   `json:"code"`
	Message   string   `json:"message"`
}

type UserAuth struct {
	UserId string `json:"userId"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Gender string `json:"gender"`
	Role   string `json:"role"`
}

type UserInfo struct {
	UserId string   `json:"userId"`
	Token  string   `json:"token"`
	Device []string `json:"device"`
	Name   string   `json:"name"`
	Avatar string   `json:"avatar"`
	Gender string   `json:"gender"`
	Role   string   `json:"role"`
}


type MessageBodySignalLogin struct {
	Token  string `json:"token"`
	Device string `json:"device"`
}

type MessageBodySignalJoinRoom struct {
	RoomId string `json:"roomId"`
}

type MessageBodySignalChangeRoom struct {
	RoomId string `json:"newRoomId"`
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
	User   UserInfo `json:"user"`
	RoomId string   `json:"roomId"`
}

type MessageBodyNoticeUnblockUser struct {
	User   UserInfo `json:"user"`
	RoomId string   `json:"roomId"`
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

type MessageBodyContentAt struct {
	AtUserId     string `json:"atUserId"`
	AtUserName   string `json:"atUserName"`
	AtUserAvatar string `json:"atUserAvatar"`
	Content      string `json:"content"`
}

type MessageBodyContentReply struct {
	ReplyUserId     string `json:"replyUserId"`
	ReplyUserName   string `json:"replyUserName"`
	ReplyUserAvatar string `json:"replyUserAvatar"`
	ReplyContent    string `json:"replyContent"`
	Content         string `json:"content"`
}

func JsonSignalLogin(any any) *MessageBodySignalLogin {
	bs, _ := json.Marshal(any)
	ret := MessageBodySignalLogin{}
	json.Unmarshal(bs, &ret)
	return &ret
}

func JsonSignalJoinRoom(any any) *MessageBodySignalJoinRoom {
	bs, _ := json.Marshal(any)
	ret := MessageBodySignalJoinRoom{}
	json.Unmarshal(bs, &ret)
	return &ret
}

func JsonSignalChangeRoom(any any) *MessageBodySignalChangeRoom {
	bs, _ := json.Marshal(any)
	ret := MessageBodySignalChangeRoom{}
	json.Unmarshal(bs, &ret)
	return &ret
}

func JsonContentText(any any) *MessageBodyContentText {
	bs, _ := json.Marshal(any)
	ret := MessageBodyContentText{}
	json.Unmarshal(bs, &ret)
	return &ret
}

func JsonContentEmoji(any any) *MessageBodyContentEmoji {
	bs, _ := json.Marshal(any)
	ret := MessageBodyContentEmoji{}
	json.Unmarshal(bs, &ret)
	return &ret
}

func JsonContentAt(any any) *MessageBodyContentAt {
	bs, _ := json.Marshal(any)
	ret := MessageBodyContentAt{}
	json.Unmarshal(bs, &ret)
	return &ret
}

func JsonContentReply(any any) *MessageBodyContentReply {
	bs, _ := json.Marshal(any)
	ret := MessageBodyContentReply{}
	json.Unmarshal(bs, &ret)
	return &ret
}
