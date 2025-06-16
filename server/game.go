package main

import (
	"image"
	"log"

	"github.com/bmcszk/gptrts/pkg/game"
	"github.com/bmcszk/gptrts/pkg/world"
)

type serverGame struct {
	*game.GameLogic
	store        game.Store
	worldService world.WorldService
	starting     map[image.Point]*game.PlayerIdType // starting point for each player, very temporary solution
}

func newServerGame(store game.Store, worldService world.WorldService) *serverGame {
	g := &serverGame{
		store:        store,
		GameLogic:    game.NewGameLogic(store),
		worldService: worldService,
		starting:     make(map[image.Point]*game.PlayerIdType),
	}
	g.starting[image.Pt(1, 1)] = nil
	g.starting[image.Pt(15, 1)] = nil
	g.starting[image.Pt(1, 15)] = nil
	g.starting[image.Pt(15, 15)] = nil

	return g
}

func (g *serverGame) HandleAction(action game.Action, dispatch game.DispatchFunc) {
	log.Printf("server handle %s", action.GetType())
	g.GameLogic.HandleAction(action, dispatch)
	switch a := action.(type) {
	case game.PlayerJoinAction:
		g.handlePlayerJoinAction(a, dispatch)
	case game.MapLoadAction:
		g.handleMapLoadAction(a, dispatch)
	}
}

func (g *serverGame) handlePlayerJoinAction(action game.PlayerJoinAction, dispatch game.DispatchFunc) {
	player := action.Payload
	id := player.Id
	_, existing := g.store.GetPlayer(id)
	g.store.StorePlayer(player)

	successAction := game.PlayerJoinSuccessAction{
		Type: game.PlayerJoinSuccessActionType,
		Payload: game.PlayerJoinSuccessPayload{
			PlayerId: player.Id,
			Units:    make([]game.Unit, 0),
			Players:  make([]game.Player, 0),
		},
	}
	for _, unit := range g.store.GetAllUnits() {
		successAction.Payload.Units = append(successAction.Payload.Units, *unit)
	}
	for _, player := range g.store.GetAllPlayers() {
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

func (g *serverGame) handleMapLoadAction(action game.MapLoadAction, dispatch game.DispatchFunc) {
	if g.isMapDataCached(action) {
		g.dispatchCachedMapData(action, dispatch)
		return
	}

	g.loadMapFromWorldService(action, dispatch)
}

func (g *serverGame) isMapDataCached(action game.MapLoadAction) bool {
	_, ok1 := g.store.GetTile(image.Pt(action.Payload.MinX, action.Payload.MinY))
	_, ok2 := g.store.GetTile(image.Pt(action.Payload.MaxX, action.Payload.MaxY))
	return ok1 && ok2
}

func (g *serverGame) dispatchCachedMapData(action game.MapLoadAction, dispatch game.DispatchFunc) {
	tiles := g.extractTilesFromStore(action)
	successAction := game.MapLoadSuccessAction{
		Type: game.MapLoadSuccessActionType,
		Payload: game.MapLoadSuccessPayload{
			WorldResponse: world.WorldResponse{Tiles: tiles},
			PlayerId:      action.Payload.PlayerId,
		},
	}
	dispatch(successAction)
}

func (g *serverGame) extractTilesFromStore(action game.MapLoadAction) []world.Tile {
	tiles := make([]world.Tile, 0)
	for x := action.Payload.MinX; x <= action.Payload.MaxX; x++ {
		for y := action.Payload.MinY; y <= action.Payload.MaxY; y++ {
			t, ok := g.store.GetTile(image.Pt(x, y))
			if !ok {
				continue
			}
			tiles = append(tiles, *t.Tile)
		}
	}
	return tiles
}

func (g *serverGame) loadMapFromWorldService(action game.MapLoadAction, dispatch game.DispatchFunc) {
	resp, err := g.worldService.Load(action.Payload.WorldRequest)
	if err != nil {
		log.Printf("error loading map: %s", err)
		return
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
