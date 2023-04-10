package main

import (
	"github.com/bmcszk/gptrts/pkg/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type Tile struct {
	*game.Tile
	x, y    int
	visible bool
}

type Map struct {
	*game.Map
	Tiles [][]*Tile
}

func NewMap(gm *game.Map) *Map {
	m := &Map{
		Map:   gm,
		Tiles: make([][]*Tile, game.MapWidth),
	}
	for x := 0; x < game.MapWidth; x++ {
		m.Tiles[x] = make([]*Tile, game.MapHeight)
		for y := 0; y < game.MapHeight; y++ {
			m.Tiles[x][y] = NewTile(x, y, gm.Tiles[x][y])
		}
	}
	return m
}

func NewTile(x, y int, t *game.Tile) *Tile {
	return &Tile{
		x:    x,
		y:    y,
		Tile: t,
	}
}

func (m *Map) Update(playerUnits []*Unit) {
	for x := 0; x < game.MapWidth; x++ {
		for y := 0; y < game.MapHeight; y++ {
			t := m.Tiles[x][y]
			t.visible = t.isVisible(playerUnits)
		}
	}
}

func (m *Map) Draw(screen *ebiten.Image, cameraX, cameraY int) {
	for x := range m.Tiles {
		for y, tile := range m.Tiles[x] {
			var img *ebiten.Image
			switch tile.Type {
			case game.Grass:
				img = grassImage
			case game.Dirt:
				img = dirtImage
			}
			op := &ebiten.DrawImageOptions{}

			if !tile.visible {
				// Create a new color matrix and set the brightness to a lower value
				cm := ebiten.ColorM{}
				cm.Scale(0.5, 0.5, 0.5, 1.0) // Make the tile darker
				op.ColorM = cm
			}

			op.GeoM.Translate(float64(x*tileSize-cameraX), float64(y*tileSize-cameraY))
			screen.DrawImage(img, op)
		}
	}
}

func (t *Tile) isVisible(playerUnits []*Unit) bool {
	p := game.NewPF(float64(t.x), float64(t.y))
	for _, u := range playerUnits {
		if p.Dist(u.Position) <= 5 {
			return true
		}
	}
	return false
}
