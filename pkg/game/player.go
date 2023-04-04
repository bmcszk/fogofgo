package game

import "github.com/google/uuid"

type PlayerIdType = uuid.UUID

type Player struct {
	Id   PlayerIdType
	Name string
}

func NewPlayer(name string) *Player {
	return &Player{
		Id:   uuid.New(),
		Name: name,
	}
}
