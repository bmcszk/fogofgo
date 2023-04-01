package main

import (
	"github.com/bmcszk/gptrts/pkg/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type Map struct {
	*game.Map
}

func NewMap() *Map {
	return &Map{
		Map: game.NewMap(),
	}
}

func (m Map) Draw(screen *ebiten.Image, cameraX, cameraY int) {
	for y := range m.Tiles {
		for x, tile := range m.Tiles[y] {
			var img *ebiten.Image
			switch tile.Type {
			case game.Grass:
				img = grassImage
			case game.Dirt:
				img = dirtImage
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*tileSize-cameraX), float64(y*tileSize-cameraY))
			screen.DrawImage(img, op)
		}
	}
}
