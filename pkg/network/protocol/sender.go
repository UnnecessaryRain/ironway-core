package protocol

// Sender interface is anything that can send a protocol message
type Sender interface {
	Send(OutgoingMessage)
}
