package main

import "github.com/hajimehoshi/ebiten/v2"

type TileType int

const (
	Grass TileType = iota
	Dirt
)

type Tile struct {
	Type TileType
}

const (
	MapWidth  = 20
	MapHeight = 20
)

type Map struct {
	Width  int
	Height int
	Tiles  [][]Tile
}

func NewMap() *Map {
	m := &Map{Tiles: make([][]Tile, MapWidth), Width: MapWidth, Height: MapHeight}
	for x := 0; x < MapWidth; x++ {
		m.Tiles[x] = make([]Tile, MapHeight)
		for y := 0; y < MapHeight; y++ {
			if x > 10 && x < 15 && y > 10 && y < 15 {
				m.Tiles[x][y] = Tile{Type: Dirt}
			} else {
				m.Tiles[x][y] = Tile{Type: Grass}
			}
		}
	}
	return m
}

func (m *Map) Draw(screen *ebiten.Image, cameraX, cameraY int) {
	for y := range m.Tiles {
		for x, tile := range m.Tiles[y] {
			var img *ebiten.Image
			switch tile.Type {
			case Grass:
				img = grassImage
			case Dirt:
				img = dirtImage
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*tileSize-cameraX), float64(y*tileSize-cameraY))
			screen.DrawImage(img, op)
		}
	}
}
