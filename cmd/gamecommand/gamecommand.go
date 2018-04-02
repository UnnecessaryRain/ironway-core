package gamecommand

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/UnnecessaryRain/ironway-core/pkg/game"
	"github.com/UnnecessaryRain/ironway-core/pkg/interpreter"
	"github.com/UnnecessaryRain/ironway-core/pkg/network/protocol"

	log "github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type stdOutSender struct{}

// Send implements the send for this stdout sender
func (s stdOutSender) Send(m protocol.Message) {
	fmt.Println(m)
}

type gameCommand struct {
	player string
}

// Configure sets up the command for gamer
func Configure(app *kingpin.Application) {
	g := &gameCommand{}
	c := app.Command("game", "starts a game").
		Action(g.run)
	// allows us to assume any user for debugging or testing
	c.Flag("assume-user", "username to assume for this game").
		Short('u').
		Required().
		StringVar(&g.player)
}

// run command for game command arg
// runs only the game accepting commands from stdin instead of websockets
func (g *gameCommand) run(c *kingpin.ParseContext) error {
	log.Infoln("Starting standalone game with player", g.player)

	// signal handling to shutdown gracefully
	sigs := make(chan os.Signal)
	stopChan := make(chan struct{})
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		close(stopChan)
		os.Exit(0)
	}()

	var outWriter stdOutSender

	gameInstance := game.NewGame()
	go gameInstance.RunForever(stopChan)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmd := interpreter.FindCommand(scanner.Text())
		gameInstance.QueueCommand(outWriter, cmd)
	}

	close(stopChan)

	return nil
}
