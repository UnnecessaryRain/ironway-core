package commands

import (
	"fmt"
	"strings"

	"github.com/UnnecessaryRain/ironway-core/pkg/game"
	"github.com/UnnecessaryRain/ironway-core/pkg/network/protocol"
)

// Debug defines a text to stdout command type
type Debug struct {
	Message string
}

// NewDebug creates a new text command
func NewDebug(message string) Debug {
	return Debug{message}
}

// Run runs this command on the game
// just print the message
func (d Debug) Run(g *game.Game) {
	fmt.Printf("Debug: %s\n", strings.TrimSpace(d.Message))
}

// Reply sends back the returned value to the sender
func (d Debug) Reply() protocol.Message {
	return protocol.Message("debug:" + d.Message)
}

// String impl method for Stringer
func (d Debug) String() string {
	return d.Message
}
