package main

import (
	"encoding/hex"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Screen struct {
	rect  image.Rectangle
	tiles [][]*Tile
}

func NewScreen() *Screen {
	return &Screen{}
}

func (m *Screen) Draw(screen *ebiten.Image, cameraX, cameraY int) {
	for _, column := range m.tiles {
		if column == nil {
			continue
		}
		for _, t := range column {
			if t == nil {
				continue
			}
			p := t.Point

			op := &ebiten.DrawImageOptions{}

			if !t.visible {
				// Create a new color matrix and set the brightness to a lower value
				cm := ebiten.ColorM{}
				cm.Scale(0.5, 0.5, 0.5, 1.0) // Make the tile darker
				op.ColorM = cm
			}

			op.GeoM.Translate(float64(p.X*tileSize-cameraX), float64(p.Y*tileSize-cameraY))
			screen.DrawImage(getBackgroundColorImage(t.BackStyleClass), op)

			tileSpriteNo := getTile(t.FrontStyleClass)
			sx := (tileSpriteNo % tileXNum) * tileSize
			sy := (tileSpriteNo / tileXNum) * tileSize

			screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
		}
	}
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
