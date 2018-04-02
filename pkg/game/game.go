package game

import (
	"github.com/UnnecessaryRain/ironway-core/pkg/network/client"
	log "github.com/sirupsen/logrus"
)

type clientCommand struct {
	client  client.Sender
	command Command
}

// Game defines the master game object and everything in the game
type Game struct {
	CommandChan chan clientCommand
}

// NewGame creates a new game object on the heap
func NewGame() *Game {
	return &Game{
		CommandChan: make(chan clientCommand, 256),
	}
}

// QueueCommand pushes the command onto the command channel for processing
func (g *Game) QueueCommand(sender client.Sender, cmd Command) {
	g.CommandChan <- clientCommand{sender, cmd}
}

// RunForever until stop channel closed
func (g *Game) RunForever(stopChan <-chan struct{}) {
	for {
		select {
		case cmd := <-g.CommandChan:
			cmd.command.Run(g)
			cmd.client.Send(cmd.command.Reply())
		case <-stopChan:
			log.Infoln("Stopping game")
			return
		}
	}
}
