package stream

import (
	"time"

	"github.com/UnnecessaryRain/ironway-core/pkg/network/protocol"
)

// Stream is a channel like stream to
type Stream struct {
	SendRate time.Duration
	Messages []protocol.OutgoingMessage
}

// NewStream creates a path to send from server to client
func NewStream(rate time.Duration) Stream {
	return Stream{
		SendRate: rate,
	}
}

// Schedule adds a new message to be sent to clients
func (p *Stream) Schedule(m protocol.OutgoingMessage) {
	p.Messages = append(p.Messages, m)
}

// Flush sends the scheduled messages
func (p *Stream) Flush(sender protocol.Sender) {
	for _, m := range p.Messages {
		sender.Send(m)
	}
	p.Messages = nil
}
