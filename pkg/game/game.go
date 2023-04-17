package game

import (
	"log"
	"sync"
)

type DispatchFunc func(Action)

type Game struct {
	Map      *Map
	Units    map[UnitIdType]*Unit
	Players  map[PlayerIdType]*Player
	Starting map[PF]*PlayerIdType
	dispatch DispatchFunc
	mux      *sync.Mutex
}

func NewGame(dispatch DispatchFunc) *Game {
	g := &Game{
		Map:      NewMap(),
		Units:    make(map[UnitIdType]*Unit),
		Players:  make(map[PlayerIdType]*Player),
		Starting: make(map[PF]*PlayerIdType),
		dispatch: dispatch,
		mux:      &sync.Mutex{},
	}

	g.Starting[PF{1, 1}] = nil
	g.Starting[PF{15, 1}] = nil
	g.Starting[PF{1, 15}] = nil
	g.Starting[PF{15, 15}] = nil
	return g
}

func (g *Game) SetMap(m *Map) {
	g.Map = m
}

func (g *Game) SetUnit(unit *Unit) {
	g.Units[unit.Id] = unit
	unit.dispatch = g.dispatch
}

func (g *Game) SetPlayer(player *Player) {
	g.Players[player.Id] = player
}

func (g *Game) HandleAction(action Action) {
	g.mux.Lock()
	defer g.mux.Unlock()

	switch a := action.(type) {
	case AddUnitAction:
		g.handleAddUnitAction(a)
	case MoveStartAction:
		g.handleMoveStartAction(a)
	case MoveStepAction:
		g.handleMoveStepAction(a)
	case MoveStopAction:
		g.handleMoveStopAction(a)
	case MapLoadSuccessAction:
		g.handleMapLoadSuccessAction(a)
	}
}

func (g *Game) handleAddUnitAction(action AddUnitAction) {
	unit := &action.Payload
	g.Units[action.Payload.Id] = unit
	if err := g.Map.PlaceUnit(unit); err != nil {
		log.Println(err)
		//dispatch error action
	}
}

func (g *Game) handleMoveStartAction(action MoveStartAction) {
	unit := g.Units[action.Payload.UnitId]

	unit.MoveTo(action.Payload.Point)
}

func (g *Game) handleMoveStepAction(action MoveStepAction) {
	//clean position
	for _, tile := range g.Map.GetTilesByUnitId(action.Payload.UnitId) {
		tile.UnitId = ZeroUnitId
	}
	unit := g.Units[action.Payload.UnitId]

	unit.Position = action.Payload.Position
	unit.Path = action.Payload.Path
	unit.Step = action.Payload.Step

	if err := g.Map.PlaceUnit(unit); err != nil {
		log.Println(err)
		//dispatch error action
	}
	//reserve next step
	if len(action.Payload.Path) > action.Payload.Step {
		nextStep := action.Payload.Path[action.Payload.Step]
		if err := g.Map.PlaceUnit(unit, nextStep); err != nil {
			g.dispatch(MoveStopAction{
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

func (g *Game) handleMoveStopAction(action MoveStopAction) {
	unit := g.Units[action.Payload]

	unit.Path = []PF{}
	unit.Step = 0

}

func (g *Game) handleMapLoadSuccessAction(action MapLoadSuccessAction) {
	for _, t := range action.Payload.Tiles {
		tile := t
		g.Map.SetTile(&tile)
	}
}
