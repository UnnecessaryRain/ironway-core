package game

import (
	"fmt"

	"github.com/UnnecessaryRain/ironway-core/pkg/network/protocol"
)

// Command interface defiens something that is a command
type Command interface {
	fmt.Stringer
	Reply() protocol.Message
	Run(*Game)
}
