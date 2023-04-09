package main

import (
	"image"
	"image/color"
	"log"
	"math"

	"github.com/bmcszk/gptrts/pkg/convert"
	"github.com/bmcszk/gptrts/pkg/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	tileSize = 32

	cameraSpeed = 2
)

type Game struct {
	*game.Game
	Map              *Map
	Units            map[game.UnitIdType]*Unit
	cameraX, cameraY int
	selectionBox     *image.Rectangle
	dispatch         game.DispatchFunc
}

func NewGame(dispatch game.DispatchFunc) *Game {
	g := &Game{
		Game:     game.NewGame(dispatch),
		Map:      NewMap(game.NewMap()),
		Units:    make(map[game.UnitIdType]*Unit),
		dispatch: dispatch,
	}
	return g
}

func (g *Game) SetMap(m *game.Map) {
	g.Game.SetMap(m)
	g.Map = NewMap(m)
}

func (g *Game) SetUnit(unit *game.Unit) {
	g.Game.SetUnit(unit)
	g.Units[unit.Id] = NewUnit(unit)
}

func (g *Game) SetPlayer(player *game.Player) {
	g.Game.SetPlayer(player)
}

func (g *Game) HandleAction(action game.Action) {
	log.Printf("client handle %s", action.GetType())
	g.Game.HandleAction(action)
	switch a := action.(type) {
	case game.AddUnitAction:
		g.handleAddUnitAction(a)
	case game.PlayerInitSuccessAction:
		g.handlePlayerInitSuccessAction(a)
	}
}

func (g *Game) handleAddUnitAction(action game.AddUnitAction) {
	g.SetUnit(&action.Payload)
}

func (g *Game) handlePlayerInitSuccessAction(action game.PlayerInitSuccessAction) {
	g.SetMap(&action.Payload.Map)
	for _, unit := range action.Payload.Units {
		g.SetUnit(&unit)
	}
	for _, player := range action.Payload.Players {
		g.SetPlayer(&player)
	}
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
			if r.Canon().Overlaps(u.GetRect()) {
				u.Selected = true
			} else if !ebiten.IsKeyPressed(ebiten.KeyShift) {
				u.Selected = false
			}
		}
	}

	// Handle right mouse button click to move selected units
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) && ebiten.IsFocused() {
		mx, my := ebiten.CursorPosition()
		worldX, worldY := g.screenToWorld(mx, my)
		tileX, tileY := worldX/tileSize, worldY/tileSize
		for _, u := range g.Units {
			if u.Selected {
				g.dispatch(game.MoveStartAction{
					Type: game.MoveStartActionType,
					Payload: game.MoveStartPayload{
						UnitId: u.Id,
						Point:  game.NewPF(float64(tileX), float64(tileY)),
					},
				})
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
