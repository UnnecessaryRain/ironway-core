package main

import (
	"os"

	"github.com/UnnecessaryRain/ironway-core/cmd/gamecmd"
	"github.com/UnnecessaryRain/ironway-core/cmd/servecmd"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("ironway", "ironway core server")
	servecmd.Configure(app)
	gamecmd.Configure(app)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
