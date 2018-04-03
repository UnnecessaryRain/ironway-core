package commands

import (
	"fmt"
	"strings"

	"github.com/UnnecessaryRain/ironway-core/pkg/mud/game"
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

// String impl method for Stringer
func (c Chat) String() string {
	return c.Message
}
