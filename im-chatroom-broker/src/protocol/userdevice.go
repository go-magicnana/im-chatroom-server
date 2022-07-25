package protocol

type UserDevice struct {
	ClientName string `json:"clientName"`
	UserId     string `json:"userId"`
	Device     string `json:"device"`
	State      string `json:"state"`
	RoomId     string `json:"roomId"`
	Broker     string `json:"broker"`
}
