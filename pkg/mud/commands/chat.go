package commands

import (
	"strings"

	"github.com/UnnecessaryRain/ironway-core/pkg/mud/game"
	log "github.com/sirupsen/logrus"
)

// Chat string to broadcast to everyone
type Chat struct {
	User      string
	Message   string
	Timestamp int64
}

// NewChat creates a new text command
func NewChat(user, message string, timestamp int64) game.Command {
	return Chat{user, message, timestamp}
}

// Run runs command on game
func (c Chat) Run(g *game.Game) {
	log.Infof("Chat: %s", strings.TrimSpace(c.Message))
	g.Chat.Post(c.User, c.Message, c.Timestamp)
}

// String impl method for Stringer
func (c Chat) String() string {
	return c.Message
}
