package game

import "fmt"

// Command interface defiens something that is a command
type Command interface {
	fmt.Stringer
	Run(*Game)
}
