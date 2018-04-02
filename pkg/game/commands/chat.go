package commands

import (
	"fmt"
	"strings"

	"github.com/UnnecessaryRain/ironway-core/pkg/game"
	"github.com/UnnecessaryRain/ironway-core/pkg/network/protocol"
)

// Chat string to broadcast to everyone
type Chat struct {
	Message string
}

// NewChat creates a new text command
func NewChat(message string) game.Command {
	return Chat{message}
}

// Run runs command on game
func (c Chat) Run(g *game.Game) {
	fmt.Printf("Chat: %s\n", strings.TrimSpace(c.Message))
}

// Reply just the chat message back to the client
func (c Chat) Reply() protocol.Message {
	return protocol.Message(c.Message)
}

// String impl method for Stringer
func (c Chat) String() string {
	return c.Message
}
