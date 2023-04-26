package main

import (
	"image"
	"log"
	"sync"

	"github.com/bmcszk/gptrts/pkg/game"
	"github.com/bmcszk/gptrts/pkg/world"
)

type Game struct {
	*game.Game
	dispatch     game.DispatchFunc
	worldService world.WorldService
	mux          *sync.Mutex
}

func NewGame(dispatch game.DispatchFunc, worldService world.WorldService) *Game {
	g := &Game{
		Game:         game.NewGame(dispatch),
		dispatch:     dispatch,
		worldService: worldService,
		mux:          &sync.Mutex{},
	}

	return g
}

func (g *Game) HandleAction(action game.Action) {
	g.mux.Lock()
	defer g.mux.Unlock()
	log.Printf("server handle %s", action.GetType())
	g.Game.HandleAction(action)
	switch a := action.(type) {
	case game.PlayerInitAction:
		g.handlePlayerInitAction(a)
	case game.MapLoadAction:
		g.handleMapLoadAction(a)
	}
}

func (g *Game) handlePlayerInitAction(action game.PlayerInitAction) {
	player := &action.Payload
	id := player.Id
	_, existing := g.Game.Players[id]
	g.Game.Players[id] = player

	successAction := game.PlayerInitSuccessAction{
		Type: game.PlayerInitSuccessActionType,
		Payload: game.PlayerInitSuccessPayload{
			PlayerId: player.Id,
			Units:    make([]game.Unit, 0),
			Players:  make([]game.Player, 0),
		},
	}
	for _, unit := range g.Game.Units {
		successAction.Payload.Units = append(successAction.Payload.Units, *unit)
	}
	for _, player := range g.Game.Players {
		successAction.Payload.Players = append(successAction.Payload.Players, *player)
	}
	g.dispatch(successAction)

	// unit spawn only for new player
	if existing {
		return
	}

	var startingP game.PF
	for sp, p := range g.Starting {
		if p == nil {
			startingP = sp
			g.Starting[sp] = &player.Id
			break
		}
	}
	unit := game.NewUnit(action.Payload.Id, player.Color, startingP, 16, 16)
	g.Units[unit.Id] = unit // should it driven by action?
	unitAction := game.SpawnUnitAction{
		Type:    game.SpawnUnitActionType,
		Payload: *unit,
	}
	g.dispatch(unitAction)
}

func (g *Game) handleMapLoadAction(action game.MapLoadAction) {

	_, ok1 := g.Map.Tiles[image.Pt(action.Payload.MinX, action.Payload.MinY)]
	_, ok2 := g.Map.Tiles[image.Pt(action.Payload.MaxX, action.Payload.MaxY)]
	if ok1 && ok2 {
		tiles := make([]world.Tile, 0)
		for x := action.Payload.MinX; x <= action.Payload.MaxX; x++ {
			for y := action.Payload.MinY; y <= action.Payload.MaxY; y++ {
				t, ok := g.Map.Tiles[image.Pt(x, y)]
				if !ok {
					continue
				}
				tiles = append(tiles, *t.Tile)
			}
		}
		successAction := game.MapLoadSuccessAction{
			Type: game.MapLoadSuccessActionType,
			Payload: game.MapLoadSuccessPayload{
				WorldResponse: world.WorldResponse{Tiles: tiles},
				PlayerId:      action.Payload.PlayerId,
			},
		}
		g.dispatch(successAction)
		return

	}

	resp, err := g.worldService.Load(action.Payload.WorldRequest)
	if err != nil {
		log.Printf("error loading map: %s", err)
		return
		// TODO: error handling
		// TODO: send error to client
	}
	successAction := game.MapLoadSuccessAction{
		Type: game.MapLoadSuccessActionType,
		Payload: game.MapLoadSuccessPayload{
			WorldResponse: *resp,
			PlayerId:      action.Payload.PlayerId,
		},
	}
	g.dispatch(successAction)
}
