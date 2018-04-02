package commands

import (
	"fmt"
	"strings"

	"github.com/UnnecessaryRain/ironway-core/pkg/game"
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
func (d Chat) Run(g *game.Game) {
	fmt.Printf("Chat: %s\n", strings.TrimSpace(d.Message))
}

// String impl method for Stringer
func (d Chat) String() string {
	return d.Message
}
