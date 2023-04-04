package main

import (
	"github.com/bmcszk/gptrts/pkg/game"
	"github.com/gorilla/websocket"
)

type Player struct {
	*game.Player
	ws *websocket.Conn
}

func NewPlayer(p *game.Player, ws *websocket.Conn) *Player {
	return &Player{
		Player: p,
		ws:     ws,
	}
}
