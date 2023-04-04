package game

import (
	"errors"

	"github.com/google/uuid"
)

type TileType int

const (
	Grass TileType = iota
	Dirt
)

type Tile struct {
	Type   TileType
	UnitId UnitIdType
}

const (
	MapWidth  = 20
	MapHeight = 20
)

type Map struct {
	Width  int
	Height int
	Tiles  [][]*Tile
}

func NewMap() *Map {
	m := &Map{Tiles: make([][]*Tile, MapWidth), Width: MapWidth, Height: MapHeight}
	for x := 0; x < MapWidth; x++ {
		m.Tiles[x] = make([]*Tile, MapHeight)
		for y := 0; y < MapHeight; y++ {
			if x > 10 && x < 15 && y > 10 && y < 15 {
				m.Tiles[x][y] = &Tile{Type: Dirt}
			} else {
				m.Tiles[x][y] = &Tile{Type: Grass}
			}
		}
	}
	return m
}

func (m *Map) GetTilesByUnitId(unitId UnitIdType) []*Tile {
	result := make([]*Tile, 0)
	for x := 0; x < MapWidth; x++ {
		for y := 0; y < MapHeight; y++ {
			tile := m.Tiles[x][y]
			if tile.UnitId == unitId {
				result = append(result, tile)
			}
		}
	}
	return result
}

func (m *Map) PlaceUnit(unit *Unit, positions ...PF) error {
	if len(positions) == 0 {
		positions = []PF{unit.Position}
	}
	for _, p := range positions {
		x, y := p.Ints()
		tile, err := m.GetTile(x, y)
		if err != nil {
			return err
		}
		//set position
		if tile.UnitId != uuid.Nil && tile.UnitId != unit.Id {
			return errors.New("position")
		}
		tile.UnitId = unit.Id
	}

	return nil
}

func (m *Map) GetTile(x, y int) (*Tile, error) {
	if x >= MapWidth || y >= MapHeight {
		return nil, errors.New("out of bounds")
	}
	return m.Tiles[x][y], nil
}
