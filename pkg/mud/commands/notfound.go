package commands

import (
	"fmt"

	"github.com/UnnecessaryRain/ironway-core/pkg/mud/game"
	"github.com/UnnecessaryRain/ironway-core/pkg/network/protocol"
	log "github.com/sirupsen/logrus"
)

// NotFound defines a canned response
type NotFound struct {
	Sender  protocol.Sender
	Message string
}

// NewNotFound creates canned response command
func NewNotFound(sender protocol.Sender, message string) game.Command {
	return NotFound{sender, message}
}

// Run command on game, sending back canned help message
func (n NotFound) Run(g *game.Game) {
	s := fmt.Sprintf("No command '%s' found. Use /help etc etc", n.Message)
	log.Infof(s)
	n.Sender.Send(protocol.OutgoingMessage{
		Frame:   protocol.LogFrame,
		Mode:    protocol.AppendMode,
		Content: fmt.Sprintf("\\red{%s}", s),
	})
}

// String impl method for Stringer
func (n NotFound) String() string {
	return n.Message
}
