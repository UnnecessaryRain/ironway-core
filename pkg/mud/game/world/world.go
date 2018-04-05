package world

// World contains the game map and entities
type World struct {
	Atlas *Atlas
}

// NewWorld creates a new world and popualtes it
func NewWorld() *World {
	return &World{
		Atlas: NewAtlas(),
	}
}
