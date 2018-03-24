package main

import (
	"os"

	"github.com/UnnecessaryRain/ironway-core/cmd/servecommand"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("ironway", "ironway core server")
	servecommand.Configure(app)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
