package main

import (
	"image"
	"sync"

	"github.com/bmcszk/gptrts/pkg/game"
	"github.com/bmcszk/gptrts/pkg/world"
)

const (
	tileSize = 16
	tileXNum = 7
)

type Tile struct {
	*game.Tile
	visible bool
}

type Map struct {
	*game.Map
	tiles *sync.Map
}

func NewMap(gm *game.Map) *Map {
	m := &Map{
		Map:   gm,
		tiles: &sync.Map{},
	}
	return m
}

func (m *Map) SetTile(tile *world.Tile) {
	p := tile.Point
	gt := m.Map.Tiles[p]
	t, ok := m.tiles.Load(p)
	if ok {
		t.(*Tile).Tile = gt
	}
	m.tiles.Store(p, &Tile{Tile: gt})
}

func (m *Map) GetTile(p image.Point) *Tile {
	t, ok := m.tiles.Load(p)
	if ok {
		return t.(*Tile)
	}
	return nil
}

func (m *Map) UpdateVisibility(unit *game.Unit) {
	m.tiles.Range(func(k any, v any) bool {
		p := k.(image.Point)
		t := v.(*Tile)

		t.visible = t.isVisible(p, unit)
		return true
	})
}

func (t *Tile) isVisible(p image.Point, playerUnits ...*game.Unit) bool {
	for _, u := range playerUnits {
		if u.Position.Dist(game.FromImagePoint(p)) <= 5 {
			return true
		}
	}
	return false
}
