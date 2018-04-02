package commands

import (
	"fmt"

	"github.com/UnnecessaryRain/ironway-core/pkg/game"
	"github.com/UnnecessaryRain/ironway-core/pkg/network/protocol"
)

// NotFound defines a canned response
type NotFound struct {
	Message string
}

// NewNotFound creates canned response command
func NewNotFound(message string) game.Command {
	return NotFound{message}
}

// Run command on game, sending back canned help message
func (d NotFound) Run(g *game.Game) {
	fmt.Printf("No command '%s' found. Use /help etc etc\n", d.Message)
}

// Reply just the not found message back to the client
func (d NotFound) Reply() protocol.Message {
	return protocol.Message(d.Message)
}

// String impl method for Stringer
func (d NotFound) String() string {
	return d.Message
}
