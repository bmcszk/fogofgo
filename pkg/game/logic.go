package game

import (
	"errors"
	"image"
	"log"
)

type DispatchFunc func(Action)

type ActionsHandler interface {
	HandleAction(Action, DispatchFunc)
}

type GameLogic struct {
	store Store
}

func NewGameLogic(store Store) *GameLogic {
	return &GameLogic{
		store: store,
	}
}

func (g *GameLogic) HandleAction(action Action, dispatch DispatchFunc) {
	switch a := action.(type) {
	case PlayerJoinSuccessAction:
		g.handlePlayerJoinSuccessAction(a, dispatch)
	case SpawnUnitAction:
		g.handleSpawnUnitAction(a, dispatch)
	case MoveStartAction:
		g.handleMoveStartAction(a)
	case MoveStepAction:
		g.handleMoveStepAction(a, dispatch)
	case MoveStopAction:
		g.handleMoveStopAction(a)
	case MapLoadSuccessAction:
		g.handleMapLoadSuccessAction(a)
	}
}

func (g *GameLogic) handlePlayerJoinSuccessAction(action PlayerJoinSuccessAction, dispatch DispatchFunc) {
	for _, u := range action.Payload.Units {
		unit := &u
		g.store.StoreUnit(unit)
		if err := g.placeUnit(unit); err != nil {
			log.Println(err)
		}
	}
	for _, p := range action.Payload.Players {
		player := p
		g.store.StorePlayer(player)
	}
}

func (g *GameLogic) handleSpawnUnitAction(action SpawnUnitAction, dispatch DispatchFunc) {
	unit := &action.Payload
	g.store.StoreUnit(unit)
	if err := g.placeUnit(unit); err != nil {
		log.Println(err)
		//dispatch error action
	}
}

func (g *GameLogic) handleMoveStartAction(action MoveStartAction) {
	unit := g.store.GetUnitById(action.Payload.UnitId)

	unit.MoveTo(action.Payload.Point)
}

func (g *GameLogic) handleMoveStepAction(action MoveStepAction, dispatch DispatchFunc) {
	//clean position
	for _, tile := range g.store.GetTilesByUnitId(action.Payload.UnitId) {
		tile.Unit = nil
	}
	unit := g.store.GetUnitById(action.Payload.UnitId)

	unit.Position = action.Payload.Position
	unit.Path = action.Payload.Path
	unit.Step = action.Payload.Step

	if err := g.placeUnit(unit); err != nil {
		log.Println(err)
		//dispatch error action
	}
	//reserve next step
	if len(action.Payload.Path) > action.Payload.Step {
		nextStep := action.Payload.Path[action.Payload.Step]
		if err := g.placeUnit(unit, nextStep); err != nil {
			dispatch(MoveStopAction{
				Type:    MoveStopActionType,
				Payload: unit.Id,
			})
			/*
				Retry moving to target in 1s
				go func(a MoveUnitAction) {
					time.Sleep(1 * time.Second)
					a.Step -= 1
					if err := g.handleAction(a); err != nil {
						log.Println(err)
					}
				}(action) */
		}
	}
}

func (g *GameLogic) handleMoveStopAction(action MoveStopAction) {
	unit := g.store.GetUnitById(action.Payload)

	unit.Path = []image.Point{}
	unit.Step = 0
}

func (g *GameLogic) handleMapLoadSuccessAction(action MapLoadSuccessAction) {
	for _, t := range action.Payload.Tiles {
		g.store.StoreTile(t)
	}
}

func (g *GameLogic) placeUnit(unit *Unit, positions ...image.Point) error {
	if len(positions) == 0 {
		positions = []image.Point{unit.Position.ImagePoint()}
	}
	for _, p := range positions {
		t, ok := g.store.GetTile(p)
		if !ok {
			t = g.store.CreateTile(p)
		}

		//set position
		if t.Unit != nil && t.Unit.Id != unit.Id {
			return errors.New("position")
		}
		t.Unit = unit
	}

	return nil
}
