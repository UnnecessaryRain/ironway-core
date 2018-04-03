package interpreter

import (
	"strings"

	"github.com/UnnecessaryRain/ironway-core/pkg/mud/commands"
	"github.com/UnnecessaryRain/ironway-core/pkg/mud/game"
	"github.com/UnnecessaryRain/ironway-core/pkg/network/client"

	log "github.com/sirupsen/logrus"
)

// commandDict list of user commands to relevant function
var commandDict = map[string]func(string) game.Command{
	"debug": commands.NewDebug,
}

// FindCommand interpretes keyword of string and return relevant Command
func FindCommand(message client.Message) game.Command {
	cmd := string(*message.Message)
	user := message.Metadata.Username
	timestamp := message.Metadata.Timestamp

	// Empty command passed somehow
	if len(cmd) == 0 {
		log.Warningln("Passed command was empty")
		return commands.NewNotFound("")
	}

	// Chat check
	if cmd[0] != '/' {
		return commands.NewChat(user, cmd, timestamp)
	}

	key := strings.Fields(cmd)[0][1:]
	if val, ok := commandDict[key]; ok {
		return val(cmd)
	}
	return commands.NewNotFound(key)
}
