package chat

import (
	"fmt"
	"time"

	"github.com/UnnecessaryRain/ironway-core/pkg/network/protocol"
)

// Chat is a chat log
type Chat struct {
	Messages []Message
}

// Format returns the message in a outgoing message format
func Format(m Message) protocol.OutgoingMessage {
	return protocol.OutgoingMessage{
		Frame:   protocol.ChatFrame,
		Content: fmt.Sprintf("\\green{%v} %v: %v", m.Timestamp.Format(time.Kitchen), m.Sender, m.Message),
		Mode:    protocol.AppendMode,
	}
}

// Flush all the messages to the sender and purge
func (c *Chat) Flush(s protocol.Sender) {
	for _, m := range c.Messages {
		s.Send(Format(m))
	}
	c.Messages = nil
}

// Post adds a new message to the chat
func (c *Chat) Post(sender, message string, timestamp int64) {
	c.Messages = append(c.Messages, Message{
		Sender:    sender,
		Message:   message,
		Timestamp: time.Unix(timestamp, 0),
	})
}

// Message is a chat message from a user
type Message struct {
	Sender    string
	Message   string
	Timestamp time.Time
}
