package game

import (
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/google/uuid"
)

type Game struct {
	Map     *Map
	Units   map[uuid.UUID]*Unit
	gameMux *sync.Mutex
	actions []any
	Dispatch func(any) error
}

func NewGame(dispatch func(any) error) *Game {
	g := &Game{
		Map: NewMap(),
		Units:   make(map[uuid.UUID]*Unit),
		gameMux: &sync.Mutex{},
	}
	/* dispatch = func(a any) error {
		if err := dispatch(a); err != nil {
			return err
		}
		if err := g.HandleAction(a); err != nil {
			return err
		}
		return nil 
	} */
	g.Dispatch = dispatch
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
		g.actions = append(g.actions, action)
	case MoveUnitAction:
		if err := g.handleMoveUnitAction(a); err != nil {
			return err
		}
		g.actions = append(g.actions, action)
	default:
		return errors.New("action not recognized")
	}
	return nil
}

func (g *Game) handleAddUnitAction(action AddUnitAction) error {
	unit := &action.Unit
	unit.dispatch = g.Dispatch
	g.Units[action.Unit.Id] = unit
	if err := g.Map.PlaceUnit(unit); err != nil {
		return err
	}
	return nil
}

func (g *Game) handleStartClientRequestAction(action StartClientRequestAction) error {
	actionJsons := []string{}
	for _, a := range g.actions {
		actionJson, err := json.Marshal(a)
		if err != nil {
			return err
		}
		actionJsons = append(actionJsons, string(actionJson))
	}

	responsAction := StartClientResponseAction {
		Type: StartClientResponseActionType,
		Actions: actionJsons,
	}
	return g.Dispatch(responsAction)
}
func (g *Game) handleStartClientResponseAction(action StartClientResponseAction) error {
	for _, actionJson := range action.Actions {
		a, err := UnmarshalAction([]byte(actionJson))
		if err != nil {
			return err
		}
		if err := g.handleAction(a); err != nil {
			return err
		}
	}
	return nil
}

func (g *Game) handleMoveUnitAction(action MoveUnitAction) error {
	//clean position
	for _, tile := range g.Map.GetTilesByUnitId(action.UnitId) {
		tile.Unit = nil
	}
	unit := g.Units[action.UnitId]

	unit.Position = action.Position
	unit.Path = action.Path
	unit.Step = action.Step

	if err := g.Map.PlaceUnit(unit); err != nil {
		return err
	}
	//reserve next step
	if len(action.Path) > action.Step {
		nextStep := action.Path[action.Step]
		if err := g.Map.PlaceUnit(unit, nextStep); err != nil {
			action.Step -= 1
			if err := g.Dispatch(action); err != nil {
				return err
			}
		}
	}
	return nil
}
