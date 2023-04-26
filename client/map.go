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
	tile := &Tile{
		Tile: &game.Tile{
			Tile: &world.Tile{
				Point: p,
			},
		},
	}
	m.tiles.Store(p, tile)
	return tile
}
