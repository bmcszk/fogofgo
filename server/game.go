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
	store        *serverStore
	worldService world.WorldService
	mux          *sync.Mutex
	starting     map[image.Point]*game.PlayerIdType
}

func NewGame(store *serverStore, worldService world.WorldService) *Game {
	g := &Game{
		store:        store,
		Game:         game.NewGame(store),
		worldService: worldService,
		mux:          &sync.Mutex{},
		starting:     make(map[image.Point]*game.PlayerIdType),
	}
	g.starting[image.Pt(1, 1)] = nil
	g.starting[image.Pt(15, 1)] = nil
	g.starting[image.Pt(1, 15)] = nil
	g.starting[image.Pt(15, 15)] = nil

	return g
}

func (g *Game) HandleAction(action game.Action, dispatch game.DispatchFunc) {
	log.Printf("server handle %s", action.GetType())
	g.Game.HandleAction(action, dispatch)
	switch a := action.(type) {
	case game.PlayerInitAction:
		g.handlePlayerInitAction(a, dispatch)
	case game.MapLoadAction:
		g.handleMapLoadAction(a, dispatch)
	}
}

func (g *Game) handlePlayerInitAction(action game.PlayerInitAction, dispatch game.DispatchFunc) {
	player := &action.Payload
	id := player.Id
	_, existing := g.store.players[id]
	g.store.players[id] = player

	successAction := game.PlayerInitSuccessAction{
		Type: game.PlayerInitSuccessActionType,
		Payload: game.PlayerInitSuccessPayload{
			PlayerId: player.Id,
			Units:    make([]game.Unit, 0),
			Players:  make([]game.Player, 0),
		},
	}
	for _, unit := range g.store.units {
		successAction.Payload.Units = append(successAction.Payload.Units, *unit)
	}
	for _, player := range g.store.players {
		successAction.Payload.Players = append(successAction.Payload.Players, *player)
	}
	dispatch(successAction)

	// unit spawn only for new player
	if existing {
		return
	}

	var startingP image.Point
	for sp, p := range g.starting {
		if p == nil {
			startingP = sp
			g.starting[sp] = &player.Id
			break
		}
	}
	unit := game.NewUnit(action.Payload.Id, player.Color, game.ToPF(startingP), 16, 16)
	unitAction := game.SpawnUnitAction{
		Type:    game.SpawnUnitActionType,
		Payload: *unit,
	}
	dispatch(unitAction)
}

func (g *Game) handleMapLoadAction(action game.MapLoadAction, dispatch game.DispatchFunc) {

	_, ok1 := g.store.tiles[image.Pt(action.Payload.MinX, action.Payload.MinY)]
	_, ok2 := g.store.tiles[image.Pt(action.Payload.MaxX, action.Payload.MaxY)]
	if ok1 && ok2 {
		tiles := make([]world.Tile, 0)
		for x := action.Payload.MinX; x <= action.Payload.MaxX; x++ {
			for y := action.Payload.MinY; y <= action.Payload.MaxY; y++ {
				t, ok := g.store.tiles[image.Pt(x, y)]
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
		dispatch(successAction)
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
	dispatch(successAction)
}
