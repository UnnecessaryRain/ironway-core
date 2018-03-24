package servecommand

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/UnnecessaryRain/ironway-core/pkg/server"
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

func (s *serveCommand) run(c *kingpin.ParseContext) error {
	log.Println("Starting server at address", s.addr)

	// signal handling to shutdown gracefully
	sigs := make(chan os.Signal)
	stopChan := make(chan struct{})
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		stopChan <- struct{}{}
		// TODO: close http server gracefully aswell
		os.Exit(0)
	}()

	server := server.NewServer(server.Options{
		Addr: s.addr,
	}, stopChan)
	server.ServeForever()

	return nil
}
