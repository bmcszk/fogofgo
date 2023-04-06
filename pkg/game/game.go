package game

import (
	"errors"
	"log"
	"sync"
)

type DispatchFunc func(any) error

type Game struct {
	Map      *Map
	Units    map[UnitIdType]*Unit
	Players  map[PlayerIdType]*Player
	Starting map[PF]*PlayerIdType
	gameMux  *sync.Mutex
	dispatch DispatchFunc
}

func NewGame(dispatch DispatchFunc) *Game {
	g := &Game{
		Map:     NewMap(),
		Units:   make(map[UnitIdType]*Unit),
		Players: make(map[PlayerIdType]*Player),
		Starting: make(map[PF]*PlayerIdType),
		gameMux: &sync.Mutex{},
	}
	g.dispatch = dispatch
	g.Starting[PF{1, 1}] = nil
	g.Starting[PF{15, 1}] = nil
	g.Starting[PF{1, 15}] = nil
	g.Starting[PF{15, 15}] = nil
	return g
}

func (g *Game) HandleAction(action any) error {
	g.gameMux.Lock()
	defer g.gameMux.Unlock()
	return g.handleAction(action)
}

func (g *Game) handleAction(action any) error {
	log.Printf("handle %+v", action)
	switch a := action.(type) {
	case StartClientRequestAction:
		if err := g.handleStartClientRequestAction(a); err != nil {
			return err
		}
	case StartClientResponseAction:
		if err := g.handleStartClientResponseAction(a); err != nil {
			return err
		}
	case AddUnitAction:
		if err := g.handleAddUnitAction(a); err != nil {
			return err
		}
	case MoveUnitAction:
		if err := g.handleMoveUnitAction(a); err != nil {
			return err
		}
	case StopUnitAction:
		if err := g.handleStopUnitAction(a); err != nil {
			return err
		}
	default:
		return errors.New("action not recognized")
	}
	return nil
}

func (g *Game) handleAddUnitAction(action AddUnitAction) error {
	unit := &action.Payload
	unit.dispatch = g.dispatch
	g.Units[action.Payload.Id] = unit
	if err := g.Map.PlaceUnit(unit); err != nil {
		return err
	}
	return nil
}

func (g *Game) handleStartClientRequestAction(action StartClientRequestAction) error {
	player := &action.Payload
	g.Players[action.Payload.Id] = player

	responsAction := StartClientResponseAction{
		Type: StartClientResponseActionType,
		Payload: StartClientResponsePayload{
			Map:     *g.Map,
			Units:   make(map[UnitIdType]Unit),
			Players: make(map[PlayerIdType]Player),
		},
	}
	for unitId, unit := range g.Units {
		responsAction.Payload.Units[unitId] = *unit
	}
	for playerId, player := range g.Players {
		responsAction.Payload.Players[playerId] = *player
	}
	if err := g.dispatch(responsAction); err != nil {
		return err
	}
	var startingP PF
	for sp, p := range g.Starting {
		if p == nil {
			startingP = sp
			g.Starting[sp] = &player.Id
			break
		}
	}
	unitAction := AddUnitAction{
		Type: AddUnitActionType,
		Payload: *NewUnit(action.Payload.Id,player.Color, startingP, 32, 32),
	}
	if err := g.dispatch(unitAction); err != nil {
		return err
	}
	return nil
}

func (g *Game) handleStartClientResponseAction(action StartClientResponseAction) error {
	g.Map = &action.Payload.Map
	for unitId, unit := range action.Payload.Units {
		gUnit := unit
		g.Units[unitId] = &gUnit
		gUnit.dispatch = g.dispatch
	}
	for playerId, player := range action.Payload.Players {
		gPlayer := player
		g.Players[playerId] = &gPlayer
	}
	return nil
}

func (g *Game) handleMoveUnitAction(action MoveUnitAction) error {
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
			if err := g.handleAction(StopUnitAction{
				Type:   StopUnitActionType,
				UnitId: unit.Id,
			}); err != nil {
				return err
			}
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

func (g *Game) handleStopUnitAction(action StopUnitAction) error {
	unit := g.Units[action.UnitId]

	unit.Path = []PF{}
	unit.Step = 0

	return nil
}
