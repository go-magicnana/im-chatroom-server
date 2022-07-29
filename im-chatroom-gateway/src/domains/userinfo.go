package domains

type UserInfo struct {
	UserId string `validate:"required" form:"userId" query:"userId" json:"userId"`
	Token  string `form:"token" query:"token" json:"token"`
	Name   string `validate:"required" form:"name" query:"name" json:"name"`
	Avatar string `validate:"required" form:"avatar" query:"avatar" json:"avatar"`
	Gender string `validate:"required" form:"gender" query:"gender" json:"gender"	`
	Role   string `form:"role" query:"role" json:"role"`
}

type BlockUser struct {
	UserId string `validate:"required" form:"userId" query:"userId" json:"userId"`
	RoomId string `validate:"required" form:"roomId" query:"roomId" json:"roomId"`
}

type Message struct{
	Command   uint16   `validate:"required" form:"command" query:"command" json:"command"`
	Target    uint32   `json:"target"`
	From      UserInfo `json:"from"`
	To        string   `json:"to"`
	Flow      uint8    `json:"flow"`
	Type      uint32   `json:"type"`
	Code      uint32   `json:"code"`
	Message   string   `json:"message"`
}
