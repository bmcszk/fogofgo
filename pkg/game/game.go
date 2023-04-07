package game

import (
	"errors"
	"sync"
)

type DispatchFunc func(Action)

type Game struct {
	Map      *Map
	Units    map[UnitIdType]*Unit
	Players  map[PlayerIdType]*Player
	Starting map[PF]*PlayerIdType
	gameMux  *sync.Mutex
	Dispatch DispatchFunc
}

func NewGame(dispatch DispatchFunc) *Game {
	g := &Game{
		Map:      NewMap(),
		Units:    make(map[UnitIdType]*Unit),
		Players:  make(map[PlayerIdType]*Player),
		Starting: make(map[PF]*PlayerIdType),
		gameMux:  &sync.Mutex{},
	}
	g.Dispatch = dispatch
	g.Starting[PF{1, 1}] = nil
	g.Starting[PF{15, 1}] = nil
	g.Starting[PF{1, 15}] = nil
	g.Starting[PF{15, 15}] = nil
	return g
}

func (g *Game) HandleAction(action Action) error {
	g.gameMux.Lock()
	defer g.gameMux.Unlock()
	return g.handleAction(action)
}

func (g *Game) handleAction(action Action) error {
	switch a := action.(type) {
	case PlayerInitAction:
		if err := g.handlePlayerInitAction(a); err != nil {
			return err
		}
	case PlayerInitSuccessAction:
		if err := g.handlePlayerInitSuccessAction(a); err != nil {
			return err
		}
	case AddUnitAction:
		if err := g.handleAddUnitAction(a); err != nil {
			return err
		}
	case MoveStartAction:
		if err := g.handleMoveStartAction(a); err != nil {
			return err
		}
	case MoveStepAction:
		if err := g.handleMoveStepAction(a); err != nil {
			return err
		}
	case MoveStopAction:
		if err := g.handleMoveStopAction(a); err != nil {
			return err
		}
	default:
		return errors.New("action not recognized")
	}
	return nil
}

func (g *Game) handleAddUnitAction(action AddUnitAction) error {
	unit := &action.Payload
	unit.dispatch = g.Dispatch
	g.Units[action.Payload.Id] = unit
	if err := g.Map.PlaceUnit(unit); err != nil {
		return err
	}
	return nil
}

func (g *Game) handleMoveStartAction(action MoveStartAction) error {
	unit := g.Units[action.Payload.UnitId]

	unit.MoveTo(action.Payload.Point)

	return nil
}

func (g *Game) handlePlayerInitAction(action PlayerInitAction) error {
	player := &action.Payload
	g.Players[action.Payload.Id] = player

	successAction := PlayerInitSuccessAction{
		Type: PlayerInitSuccessActionType,
		Payload: PlayerInitSuccessPayload{
			PlayerId: player.Id,
			Map:      *g.Map,
			Units:    make(map[UnitIdType]Unit),
			Players:  make(map[PlayerIdType]Player),
		},
	}
	for unitId, unit := range g.Units {
		successAction.Payload.Units[unitId] = *unit
	}
	for playerId, player := range g.Players {
		successAction.Payload.Players[playerId] = *player
	}
	g.Dispatch(successAction)

	var startingP PF
	for sp, p := range g.Starting {
		if p == nil {
			startingP = sp
			g.Starting[sp] = &player.Id
			break
		}
	}
	unit := NewUnit(action.Payload.Id, player.Color, startingP, 32, 32)
	g.Units[unit.Id] = unit // should it driven by action?
	unitAction := AddUnitAction{
		Type:    AddUnitActionType,
		Payload: *unit,
	}
	g.Dispatch(unitAction)

	return nil
}

func (g *Game) handlePlayerInitSuccessAction(action PlayerInitSuccessAction) error {
	g.Map = &action.Payload.Map
	for unitId, unit := range action.Payload.Units {
		gUnit := unit
		g.Units[unitId] = &gUnit
		gUnit.dispatch = g.Dispatch
	}
	for playerId, player := range action.Payload.Players {
		gPlayer := player
		g.Players[playerId] = &gPlayer
	}
	return nil
}

func (g *Game) handleMoveStepAction(action MoveStepAction) error {
	//clean position
	for _, tile := range g.Map.GetTilesByUnitId(action.Payload.UnitId) {
		tile.UnitId = ZeroUnitId
	}
	unit := g.Units[action.Payload.UnitId]

	unit.Position = action.Payload.Position
	unit.Path = action.Payload.Path
	unit.Step = action.Payload.Step

	if err := g.Map.PlaceUnit(unit); err != nil {
		return err
	}
	//reserve next step
	if len(action.Payload.Path) > action.Payload.Step {
		nextStep := action.Payload.Path[action.Payload.Step]
		if err := g.Map.PlaceUnit(unit, nextStep); err != nil {
			g.Dispatch(MoveStopAction{
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
	return nil
}

func (g *Game) handleMoveStopAction(action MoveStopAction) error {
	unit := g.Units[action.Payload]

	unit.Path = []PF{}
	unit.Step = 0

	return nil
}
