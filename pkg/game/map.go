package game

import (
	"errors"
	"image"

	"github.com/bmcszk/gptrts/pkg/world"
)

type Tile struct {
	*world.Tile
	UnitId UnitIdType
}

type Map struct {
	Tiles map[image.Point]*Tile
}

func NewMap() *Map {
	m := &Map{Tiles: make(map[image.Point]*Tile, 2000)}
	return m
}

func (m *Map) GetTilesByUnitId(unitId UnitIdType) []*Tile {
	result := make([]*Tile, 0)
	for _, t := range m.Tiles {
		if t.UnitId == unitId {
			result = append(result, t)
		}
	}
	return result
}

func (m *Map) PlaceUnit(unit *Unit, positions ...PF) error {
	if len(positions) == 0 {
		positions = []PF{unit.Position}
	}
	for _, pf := range positions {
		p := pf.ImagePoint()
		t, ok := m.Tiles[p]
		if !ok {
			m.Tiles[p] = &Tile{UnitId: unit.Id}
			return nil
		}

		//set position
		if t.UnitId != ZeroUnitId && t.UnitId != unit.Id {
			return errors.New("position")
		}
		t.UnitId = unit.Id
	}

	return nil
}

func (m *Map) SetTile(tile *world.Tile) {
	p := tile.Point
	t, ok := m.Tiles[p]
	if ok {
		t.Tile = tile
	}
	m.Tiles[p] = &Tile{Tile: tile}
}
