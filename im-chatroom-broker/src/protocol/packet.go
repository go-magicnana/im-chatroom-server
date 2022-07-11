package protocol

const (
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
	TypeSignalConnect    = 2102
	TypeSignalDisconnect = 2103
	TypeSignalJoinRoom   = 2104
	TypeSignalLeaveRoom  = 2105

	TypeNoticeJoinRoom    = 3101
	TypeNoticeLeaveRoom   = 3102
	TypeNoticeBlockUser   = 3103
	TypeNoticeUnblockUser = 3104
	TypeNoticeCloseRoom   = 3105
	TypeNoticeBlockRoom   = 3105
	TypeNoticeUnblockRoom = 3105

	TypeContentText  = 4105
	TypeContentEmoji = 4105
	TypeContentReply = 4105

	TypeGiftNone = 5101

	TypeGoodsNone = 6101
)

var ErrorOK = ImError{Code: 200, Message: "OK"}
var ErrorDefault = ImError{Code: 1001, Message: "Server Error"}
var ErrorDefaultKey = ImError{1002,"Command Not Allow"}

type ImError struct {
	Code    uint32
	Message string
}

type Packet struct {
	MessageId string
	Version   uint8
	Command   uint16
	Message   Message
}

type Message struct {
	Header MessageHeader
	Body   any
}

func NewResponseOK(in *Packet, body any) *Packet {

	header := MessageHeader{
		Target:  in.Message.Header.Target,
		From:    in.Message.Header.From,
		To:      in.Message.Header.To,
		Type:    in.Message.Header.Type,
		Flow:    FlowDown,
		Code:    ErrorOK.Code,
		Message: ErrorOK.Message,
	}

	message := Message{
		Header: header,
		Body:   body,
	}

	return &Packet{
		MessageId: in.MessageId,
		Version:   in.Version,
		Command:   in.Command,
		Message:   message,
	}
}

func NewResponseError(in *Packet, error *ImError) *Packet {
	header := MessageHeader{
		Target:  in.Message.Header.Target,
		From:    in.Message.Header.From,
		To:      in.Message.Header.To,
		Type:    in.Message.Header.Type,
		Flow:    FlowDown,
		Code:    error.Code,
		Message: error.Message,
	}

	message := Message{
		Header: header,
	}

	return &Packet{
		MessageId: in.MessageId,
		Version:   in.Version,
		Command:   in.Command,
		Message:   message,
	}
}

type MessageHeader struct {
	Target  uint32 `json:"target"`
	From    User   `json:"from"`
	To      User   `json:"to"`
	Flow    uint8  `json:"flow"`
	Type    uint32 `json:"type"`
	Code    uint32 `json:"code"`
	Message string `json:"message"`
}

type User struct {
	UserId   uint64 `json:"userId"`
	UserName string `json:"userName"`
	Avatar   string `json:"avatar"`
	Role     uint32 `json:"role"`
}

type MessageBodySignalPing struct {
	Current uint64 `json:"current"`
}

type MessageBodySignalConnect struct {
	UserId string `json:"userId"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Role   uint8  `json:"role"`
}

type MessageBodySignalDisconnect struct {
	UserId string `json:"userId"`
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
