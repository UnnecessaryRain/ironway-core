package game

import (
	log "github.com/sirupsen/logrus"
)

// Game defines the master game object and everything in the game
type Game struct {
	CommandChan chan Command
}

// NewGame creates a new game object on the heap
func NewGame() *Game {
	return &Game{
		CommandChan: make(chan Command, 256),
	}
}

// QueueCommand pushes the command onto the command channel for processing
func (g *Game) QueueCommand(cmd Command) {
	g.CommandChan <- cmd
}

// RunForever until stop channel closed
func (g *Game) RunForever(stopChan <-chan struct{}) {
	for {
		select {
		case cmd := <-g.CommandChan:
			cmd.Run(g)
		case <-stopChan:
			log.Infoln("Stopping game")
			return
		}
	}
}
