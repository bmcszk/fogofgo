package game

import "github.com/google/uuid"

type Game struct {
	Map              *Map
	Units            map[uuid.UUID]*Unit
}

func NewGame() *Game {
	return &Game{
		Units: make(map[uuid.UUID]*Unit),
	}
}
