package main

import (
	"image"
	"image/color"
	"math"

	"github.com/bmcszk/gptrts/pkg/convert"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	tileSize = 32

	cameraSpeed = 2
)

type Game struct {
	Map              *Map
	Units            []*Unit
	cameraX, cameraY int
	selectionBox     *image.Rectangle
}

func (g *Game) Init() {
	// Initialize the map tiles
	g.Map = NewMap()
	unit1 := &Unit{
		X: 3, Y: 3, Color: color.RGBA{255, 0, 0, 255}, Width: 32, Height: 32,
	}

	unit2 := &Unit{
		X: 20, Y: 10, Color: color.RGBA{0, 0, 255, 255}, Width: 32, Height: 32,
	}

	g.Units = append(g.Units, unit1, unit2)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Calculate the desired screen size based on the size of the map
	screenWidth = len(g.Map.Tiles[0]) * tileSize
	screenHeight = len(g.Map.Tiles) * tileSize

	// Scale the screen if it is too large to fit
	if screenWidth > outsideWidth || screenHeight > outsideHeight {
		scale := math.Min(float64(outsideWidth)/float64(screenWidth), float64(outsideHeight)/float64(screenHeight))
		screenWidth = int(float64(screenWidth) * scale)
		screenHeight = int(float64(screenHeight) * scale)
	}

	return screenWidth, screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the map
	g.Map.Draw(screen, g.cameraX, g.cameraY)

	// Draw units
	for _, unit := range g.Units {
		unit.Draw(screen, g.cameraX, g.cameraY)
	}

	// Draw the selection box
	if g.selectionBox != nil {
		r := *g.selectionBox
		x1, y1 := g.worldToScreen(r.Min.X, r.Min.Y)
		x2, y2 := g.worldToScreen(r.Max.X, r.Max.Y)

		col := color.RGBA{0, 255, 0, 128}
		ebitenutil.DrawRect(screen, float64(x1), float64(y1), float64(x2-x1), float64(y2-y1), col)
	}
}

func (g *Game) Update() error {
	// Move camera with arrow keys
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.cameraX -= cameraSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.cameraX += cameraSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.cameraY -= cameraSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.cameraY += cameraSpeed
	}

	// Handle left mouse button click to select units
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && ebiten.IsFocused() {
		mx, my := ebiten.CursorPosition()
		worldX, worldY := g.screenToWorld(mx, my)

		if g.selectionBox == nil {
			g.selectionBox = convert.ToPointer(image.Rect(worldX, worldY, worldX+1, worldY+1))
		} else {
			g.selectionBox.Max = image.Pt(worldX+1, worldY+1)
		}
	} else {
		g.selectionBox = nil
	}

	if g.selectionBox != nil {
		r := *g.selectionBox
		for _, u := range g.Units {
			unitRect := image.Rect(int(u.X), int(u.Y), int(u.X)+int(u.Width), int(u.Y)+int(u.Height))
			if r.Canon().Overlaps(unitRect) {
				u.Selected = true
			} else {
				u.Selected = false
			}
		}
	}

	// Handle right mouse button click to move selected units
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) && ebiten.IsFocused() {
		mx, my := ebiten.CursorPosition()
		worldX, worldY := g.screenToWorld(mx, my)
		for _, u := range g.Units {
			if u.Selected {
				u.MoveTo(worldX, worldY)
			}
		}
	}

	for _, u := range g.Units {
		err := u.Update()
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Game) screenToWorld(screenX, screenY int) (int, int) {
	worldX := screenX + g.cameraX
	worldY := screenY + g.cameraY
	return worldX, worldY
}

func (g *Game) worldToScreen(worldX, worldY int) (int, int) {
	screenX := worldX - g.cameraX
	screenY := worldY - g.cameraY
	return screenX, screenY
}
