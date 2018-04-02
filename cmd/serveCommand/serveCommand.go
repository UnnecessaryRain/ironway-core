package servecommand

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/UnnecessaryRain/ironway-core/pkg/game/commands"

	"github.com/UnnecessaryRain/ironway-core/pkg/game"
	"github.com/UnnecessaryRain/ironway-core/pkg/network/client"

	"github.com/UnnecessaryRain/ironway-core/pkg/network/server"
	log "github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type serveCommand struct {
	addr string
}

// Configure sets up the command for server
func Configure(app *kingpin.Application) {
	s := &serveCommand{}
	c := app.Command("serve", "starts a server").
		Action(s.run)

	c.Flag("addr", "address:port used to bind server").
		Default(":8080").
		StringVar(&s.addr)
}

// run the serve command code
// creates the game and the server and pushes messages from
// the server to the client
func (s *serveCommand) run(c *kingpin.ParseContext) error {
	log.Infoln("Starting server at address", s.addr)

	// signal handling to shutdown gracefully
	sigs := make(chan os.Signal)
	stopChan := make(chan struct{})
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		stopChan <- struct{}{}
		// TODO(#4): close http server gracefully aswell
		os.Exit(0)
	}()

	gameInstance := game.NewGame()
	go gameInstance.RunForever(stopChan)

	server := server.NewServer(server.Options{
		Addr: s.addr,
	})
	server.OnMessage(func(m client.Message) {
		gameInstance.QueueCommand(m.Client, commands.NewDebug(string(*m.Message)))
	})
	server.ServeForever(stopChan)

	return nil
}
