package game

import "github.com/google/uuid"

type PlayerIdType struct {
	uuid.UUID
}

type Player struct {
	Id   PlayerIdType
	Name string
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
