package world

type World struct {
	Atlas *Atlas
}

func NewWorld() *World {
	return &World{
		Atlas: NewAtlas(),
	}
}
