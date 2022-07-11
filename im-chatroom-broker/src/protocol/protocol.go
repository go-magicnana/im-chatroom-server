package protocol

type abstractProtocol interface {
	Enpack(message []byte) interface{}
	Depack(message interface{}) []byte
}
