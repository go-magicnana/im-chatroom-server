package protocol

const (
	CmdWrite = "write"
	CmdQuit  = "quit"
)

type InnerPacket struct {
	Cmd        string
	Packet     *Packet
	ClientName string
	UserId     string
	RoomId     string
}

func NewResponse(p *Packet) *InnerPacket {
	return &InnerPacket{
		Cmd:    CmdWrite,
		Packet: p,
	}
}

func NewQuit() *InnerPacket {
	return &InnerPacket{
		Cmd:    CmdQuit,
		Packet: nil,
	}
}
