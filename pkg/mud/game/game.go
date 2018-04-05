package game

import (
	"time"

	"github.com/UnnecessaryRain/ironway-core/pkg/network/protocol"
	"github.com/UnnecessaryRain/ironway-core/pkg/network/stream"
	log "github.com/sirupsen/logrus"
)

type clientCommand struct {
	client  protocol.Sender
	command Command
}

// Game defines the master game object and everything in the game
type Game struct {
	Chat        stream.Stream
	CommandChan chan clientCommand
	everyone    protocol.Sender
}

// NewGame creates a new game object on the heap
func NewGame(broadcaster protocol.Sender) *Game {
	return &Game{
		Chat:        stream.NewStream(100),
		CommandChan: make(chan clientCommand, 256),
		everyone:    broadcaster,
	}
}

// QueueCommand pushes the command onto the command channel for processing
func (g *Game) QueueCommand(sender protocol.Sender, cmd Command) {
	g.CommandChan <- clientCommand{sender, cmd}
}

// RunForever until stop channel closed
func (g *Game) RunForever(stopChan <-chan struct{}) {
	chatTicker := time.NewTicker(g.Chat.SendRate * time.Millisecond)
	for {
		select {
		case <-chatTicker.C:
			g.Chat.Flush(g.everyone)
		case cmd := <-g.CommandChan:
			cmd.command.Run(g)
		case <-stopChan:
			log.Infoln("Stopping game")
			chatTicker.Stop()
			return
		}
	}
}
