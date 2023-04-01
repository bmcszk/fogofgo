package game

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type Game struct {
	Map     *Map
	Units   map[uuid.UUID]*Unit
	PxU     map[PF]uuid.UUID
	unitMux *sync.Mutex
}

func NewGame() *Game {
	return &Game{
		Units:   make(map[uuid.UUID]*Unit),
		PxU:     make(map[PF]uuid.UUID),
		unitMux: &sync.Mutex{},
	}
}

func (g *Game) AddUnit(unit *Unit) error {
	g.unitMux.Lock()
	defer g.unitMux.Unlock()
	g.Units[unit.Id] = unit
	//clean position
	for p, id := range g.PxU {
		if id == unit.Id {
			delete(g.PxU, p)
		}
	}
	//set position
	if id, exists := g.PxU[unit.Position]; exists && id != unit.Id {
		return errors.New("position")
	}
	g.PxU[unit.Position] = unit.Id
	//reserve next step
	if len(unit.Path) > unit.Step {
		nextStep := unit.Path[unit.Step]
		if id, exists := g.PxU[nextStep]; exists && id != unit.Id {
			return errors.New("next step")
		}
		g.PxU[nextStep] = unit.Id
	}
	return nil
}
