package main

import (
	"encoding/hex"
	"image"
	"image/color"
	"log"
	"sync"

	"github.com/bmcszk/gptrts/pkg/game"
	"github.com/bmcszk/gptrts/pkg/world"
	"github.com/hajimehoshi/ebiten/v2"
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
	tiles                  *sync.Map
	minX, minY, maxX, maxY int
}

func NewMap(gm *game.Map) *Map {
	m := &Map{
		Map:   gm,
		tiles: &sync.Map{},
	}
	return m
}

func (m *Map) SetTile(tile *world.Tile) {
	p := game.NewPF(float64(tile.Point.X), float64(tile.Point.Y))
	gt := m.Map.Tiles[p]
	t, ok := m.tiles.Load(p)
	if ok {
		t.(*Tile).Tile = gt
	}
	m.tiles.Store(p, &Tile{Tile: gt})
}

func (g *Game) SetVisibleTiles(minX, minY, maxX, maxY int) {
	m := g.Map
	if minX < m.minX {
		action := game.MapLoadAction{
			Type: game.MapLoadActionType,
			Payload: game.MapLoadPayload{
				WorldRequest: world.WorldRequest{
					MinX: minX,
					MinY: minY,
					MaxX: m.minX,
					MaxY: m.maxY,
				},
				PlayerId: g.PlayerId,
			},
		}
		g.dispatch(action)
		m.minX = minX
	}
	if minY < m.minY {
		action := game.MapLoadAction{
			Type: game.MapLoadActionType,
			Payload: game.MapLoadPayload{
				WorldRequest: world.WorldRequest{
					MinX: minX,
					MinY: minY,
					MaxX: m.maxX,
					MaxY: m.minY,
				},
				PlayerId: g.PlayerId,
			},
		}
		g.dispatch(action)
		m.minY = minY
	}
	if maxX > m.maxX {
		action := game.MapLoadAction{
			Type: game.MapLoadActionType,
			Payload: game.MapLoadPayload{
				WorldRequest: world.WorldRequest{
					MinX: m.minX,
					MinY: m.maxY,
					MaxX: maxX,
					MaxY: maxY,
				},
				PlayerId: g.PlayerId,
			},
		}
		g.dispatch(action)
		m.maxX = maxX
	}
	if maxY > m.maxY {
		action := game.MapLoadAction{
			Type: game.MapLoadActionType,
			Payload: game.MapLoadPayload{
				WorldRequest: world.WorldRequest{
					MinX: m.maxY,
					MinY: m.minY,
					MaxX: maxX,
					MaxY: maxY,
				},
				PlayerId: g.PlayerId,
			},
		}
		g.dispatch(action)
		m.maxY = maxY
	}
}

func (m *Map) UpdateVisibility(unit *game.Unit) {
	m.tiles.Range(func(k any, v any) bool {
		p := k.(game.PF)
		t := v.(*Tile)

		t.visible = t.isVisible(p, unit)
		return true
	})
}

func (m *Map) Draw(screen *ebiten.Image, cameraX, cameraY int) {
	m.tiles.Range(func(k any, v any) bool {
		p := k.(game.PF)
		t := v.(*Tile)

		op := &ebiten.DrawImageOptions{}

		if !t.visible {
			// Create a new color matrix and set the brightness to a lower value
			cm := ebiten.ColorM{}
			cm.Scale(0.5, 0.5, 0.5, 1.0) // Make the tile darker
			op.ColorM = cm
		}

		op.GeoM.Translate(p.X*tileSize-float64(cameraX), p.Y*tileSize-float64(cameraY))
		screen.DrawImage(getBackgroundColorImage(t.BackStyleClass), op)

		tileSpriteNo := getTile(t.FrontStyleClass)
		sx := (tileSpriteNo % tileXNum) * tileSize
		sy := (tileSpriteNo / tileXNum) * tileSize

		screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)

		return true
	})
}

func (t *Tile) isVisible(p game.PF, playerUnits ...*game.Unit) bool {
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
