package main

import (
	"encoding/hex"
	"image"
	"image/color"
	"log"

	"github.com/bmcszk/gptrts/pkg/game"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	tileSize = 16
	tileXNum = 7
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
			op := &ebiten.DrawImageOptions{}

			if !tile.visible {
				// Create a new color matrix and set the brightness to a lower value
				cm := ebiten.ColorM{}
				cm.Scale(0.5, 0.5, 0.5, 1.0) // Make the tile darker
				op.ColorM = cm
			}

			op.GeoM.Translate(float64(x*tileSize-cameraX), float64(y*tileSize-cameraY))
			screen.DrawImage(getBackgroundColorImage(tile.BackStyleClass), op)

			t := getTile(tile.FrontStyleClass)
			sx := (t % tileXNum) * tileSize
			sy := (t / tileXNum) * tileSize

			screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
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

func getBackgroundColorImage(className string) *ebiten.Image {
	img, exists := backgroundImages[className]
	if exists {
		return img
	}
	switch className {
	case "grass":
		img = ebiten.NewImage(tileSize, tileSize)
		img.Fill(getColorFromHex("74CF45FF"))
	case "water":
		img = ebiten.NewImage(tileSize, tileSize)
		img.Fill(getColorFromHex("68A8C8FF"))
	case "sand":
		img = ebiten.NewImage(tileSize, tileSize)
		img.Fill(getColorFromHex("BA936BFF"))
	default:
		img = ebiten.NewImage(tileSize, tileSize)
		img.Fill(getColorFromHex("74CF45FF"))
	}
	backgroundImages[className] = img
	return img
}

func getColorFromHex(colorStr string) color.Color {
	b, err := hex.DecodeString(colorStr)
	if err != nil {
		log.Fatal(err)
	}

	return color.RGBA{b[0], b[1], b[2], b[3]}
}

func getTile(className string) int {
	switch className {
	case "plain1":
		return 0
	case "plain2":
		return 1
	case "plain3":
		return 2
	case "forest1":
		return 3
	case "forest2":
		return 4
	case "forest3":
		return 5
	case "sea1", "sea2", "sea3", "river1", "river2", "river3":
		return 99
	case "mountain1":
		return 10
	case "mountain2":
		return 11
	case "mountain3":
		return 12
	case "hill1", "hill2", "hill3":
		return 13
	case "lake1":
		return 18
	case "lake2":
		return 19
	case "lake3":
		return 20
	case "sand1", "sand2", "sand3":
		return 78
	default:
		return 0
	}
}