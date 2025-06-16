package main

import (
	"encoding/hex"
	"image"
	"image/color"
	"log"

	"github.com/bmcszk/fogofgo/pkg/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	tileSize       = 16
	tileSpriteSize = 16
	tileSpriteXNum = 7
	selectedBorder = 2
)

type screen struct {
	rect  image.Rectangle
	tiles map[image.Point]*game.Tile
	units map[*game.Unit]bool
}

var emptyScreen = screen{
	rect:  image.Rectangle{},
	tiles: make(map[image.Point]*game.Tile),
	units: make(map[*game.Unit]bool),
}

func newScreen(rect image.Rectangle, tiles map[image.Point]*game.Tile) *screen {
	return &screen{
		rect:  rect,
		tiles: tiles,
		units: make(map[*game.Unit]bool, 0),
	}
}

func (s *screen) is(rect image.Rectangle) bool {
	return s.rect.Eq(rect)
}

func (s *screen) draw(enScreen *ebiten.Image, cameraX, cameraY int) {
	s.drawTiles(enScreen, cameraX, cameraY)
	s.drawVisibleUnits(enScreen, cameraX, cameraY)
}

func (s *screen) drawTiles(enScreen *ebiten.Image, cameraX, cameraY int) {
	for _, t := range s.tiles {
		if t != nil {
			drawTile(t, enScreen, cameraX, cameraY)
			s.trackUnitVisibility(t)
		}
	}
}

func (s *screen) trackUnitVisibility(t *game.Tile) {
	if t.Unit != nil {
		s.units[t.Unit] = t.Visible
	}
}

func (s *screen) drawVisibleUnits(enScreen *ebiten.Image, cameraX, cameraY int) {
	for u, visible := range s.units {
		if visible {
			drawUnit(u, enScreen, cameraX, cameraY)
		}
	}
}

func drawTile(t *game.Tile, enScreen *ebiten.Image, cameraX, cameraY int) {
	p := t.Point

	op := &ebiten.DrawImageOptions{}

	if !t.Visible {
		// Use ColorScale for dimming instead of deprecated ColorM
		op.ColorScale.Scale(0.5, 0.5, 0.5, 1.0) // Make the tile darker
	}

	op.GeoM.Translate(float64(p.X*tileSize-cameraX), float64(p.Y*tileSize-cameraY))
	enScreen.DrawImage(getBackgroundColorImage(t.BackStyleClass), op)

	tileSpriteNo := getTile(t.FrontStyleClass)
	sx := (tileSpriteNo % tileSpriteXNum) * tileSpriteSize
	sy := (tileSpriteNo / tileSpriteXNum) * tileSpriteSize

	subImage := tilesImage.SubImage(image.Rect(sx, sy, sx+tileSpriteSize, sy+tileSpriteSize))
	enScreen.DrawImage(subImage.(*ebiten.Image), op)
}

func drawUnit(u *game.Unit, enScreen *ebiten.Image, cameraX, cameraY int) {
	screenPosition := u.Position.Mul(tileSize)
	x := screenPosition.X - float64(cameraX)
	y := screenPosition.Y - float64(cameraY)

	if u.Selected {
		col := color.RGBA{0, 255, 0, 255}
		borderWidth := float32(u.Size.X + (selectedBorder * 2))
		borderHeight := float32(u.Size.Y + (selectedBorder * 2))
		vector.DrawFilledRect(enScreen, float32(x-selectedBorder), float32(y-selectedBorder),
			borderWidth, borderHeight, col, false)
	}

	vector.DrawFilledRect(enScreen, float32(x), float32(y), float32(u.Size.X), float32(u.Size.Y), u.Color, false)
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
		log.Printf("Error decoding hex color %s: %v", colorStr, err)
		return color.RGBA{0, 0, 0, 0}
	}

	return color.RGBA{b[0], b[1], b[2], b[3]}
}

var tileMap = map[string]int{
	"plain1":    0,
	"plain2":    1,
	"plain3":    2,
	"forest1":   3,
	"forest2":   4,
	"forest3":   5,
	"sea1":      99,
	"sea2":      99,
	"sea3":      99,
	"river1":    99,
	"river2":    99,
	"river3":    99,
	"mountain1": 10,
	"mountain2": 11,
	"mountain3": 12,
	"hill1":     13,
	"hill2":     13,
	"hill3":     13,
	"lake1":     18,
	"lake2":     19,
	"lake3":     20,
	"sand1":     78,
	"sand2":     78,
	"sand3":     78,
}

func getTile(className string) int {
	if tileNum, exists := tileMap[className]; exists {
		return tileNum
	}
	return 0
}
