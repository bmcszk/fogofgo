package main

import (
	"sync"

	"github.com/bmcszk/gptrts/pkg/game"
	"github.com/gorilla/websocket"
)

type Game struct {
	*game.Game
	Players map[game.PlayerIdType]*Player
	gameMux *sync.Mutex
}

func NewGame(dispatch game.DispatchFunc) *Game {
	return &Game{
		Game:    game.NewGame(dispatch),
		Players: make(map[game.PlayerIdType]*Player),
		gameMux: &sync.Mutex{},
	}
}

func (g *Game) HandleAction(action game.Action, ws *websocket.Conn) error {
	g.gameMux.Lock()
	defer g.gameMux.Unlock()
	if err := g.handleAction(action, ws); err != nil {
		return err
	}
	if err := g.Game.HandleAction(action); err != nil {
		return err
	}
	return nil
}

func (g *Game) handleAction(action any, ws *websocket.Conn) error {
	switch a := action.(type) {
	case game.PlayerInitAction:
		if err := g.handleStartClientRequestAction(a, ws); err != nil {
			return err
		}
	}
	return nil
}

func (g *Game) handleStartClientRequestAction(action game.PlayerInitAction, ws *websocket.Conn) error {
	g.Players[action.Payload.Id] = NewPlayer(&action.Payload, ws)
	return nil
}
