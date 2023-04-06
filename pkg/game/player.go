package game

import (
	"image/color"

	"github.com/google/uuid"
)

type PlayerIdType struct {
	uuid.UUID
}

type Player struct {
	Id   PlayerIdType
	Name string
	Color color.RGBA
	Start PF
}

func NewPlayer(name string) *Player {
	return &Player{
		Id:   NewPlayerId(),
		Name: name,
	}
}

func NewPlayerId() PlayerIdType {
	return PlayerIdType{uuid.New()}
}
