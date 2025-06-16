package main

import (
	"image"
	"image/color"
	"log"

	"github.com/bmcszk/fogofgo/pkg/convert"
	"github.com/bmcszk/fogofgo/pkg/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	cameraSpeed = 2
)

type clientGame struct {
	*game.GameLogic
	store            game.Store
	playerId         game.PlayerIdType
	cameraX, cameraY int
	centerX, centerY int
	selectionBox     *image.Rectangle
	enDispatch       game.DispatchFunc
	screen           *screen
}

func newClientGame(playerId game.PlayerIdType, store game.Store, enDispatch game.DispatchFunc) *clientGame {
	g := game.NewGameLogic(store)
	cg := &clientGame{
		store:      store,
		playerId:   playerId,
		GameLogic:  g,
		enDispatch: enDispatch,
		screen:     &emptyScreen,
	}

	return cg
}

func (g *clientGame) HandleAction(action game.Action, dispatch game.DispatchFunc) {
	log.Printf("client handle %s", action.GetType())
	g.GameLogic.HandleAction(action, dispatch)
	switch action.(type) {
	case game.SpawnUnitAction, game.MoveStepAction, game.PlayerJoinSuccessAction, game.MapLoadSuccessAction:
		g.updateVisibility()
	}
}

func (g *clientGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	/* // Calculate the desired screen size based on the size of the map
	sw := len(g.Map.Tiles[0]) * tileSize
	sh := len(g.Map.Tiles) * tileSize

	// Scale the screen if it is too large to fit
	if sw > outsideWidth || sh > outsideHeight {
		scale := math.Min(float64(outsideWidth)/float64(sw), float64(outsideHeight)/float64(sh))
		sw = int(float64(sw) * scale)
		sh = int(float64(sh) * scale)
	}

	return sw, sh */

	// g.centerX = -outsideWidth /2
	// g.centerY = -outsideHeight /2

	minX, minY := g.screenToWorldTiles(0, 0)
	maxX, maxY := g.screenToWorldTiles(outsideWidth, outsideHeight)
	rect := image.Rect(minX, minY, maxX, maxY)

	// If the map is not loaded, load it
	if !g.screen.is(rect) {
		g.queueMapLoadActions(rect)
		g.screen = newScreen(rect, g.store.GetTilesByRect(rect))
		g.updateVisibility()
	}

	return outsideWidth, outsideHeight
}

func (g *clientGame) Draw(enScreen *ebiten.Image) {
	// Draw the map
	g.screen.draw(enScreen, g.centerX+g.cameraX, g.centerY+g.cameraY)

	// Draw the selection box
	if g.selectionBox != nil {
		r := *g.selectionBox
		x1, y1 := g.worldToScreen(r.Min.X, r.Min.Y)
		x2, y2 := g.worldToScreen(r.Max.X, r.Max.Y)

		col := color.RGBA{0, 255, 0, 128}
		vector.DrawFilledRect(enScreen, float32(x1), float32(y1), float32(x2-x1), float32(y2-y1), col, false)
	}
}

func (g *clientGame) Update() error {
	g.handleCameraMovement()
	g.handleUnitSelection()
	g.handleUnitMovement()
	g.updateUnits()
	return nil
}

func (g *clientGame) handleCameraMovement() {
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
}

func (g *clientGame) handleUnitSelection() {
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
		g.updateUnitSelectionFromBox()
	}
}

func (g *clientGame) updateUnitSelectionFromBox() {
	r := *g.selectionBox
	for _, u := range g.store.GetAllUnits() {
		if r.Canon().Overlaps(getRect(u)) {
			u.Selected = true
		} else if !ebiten.IsKeyPressed(ebiten.KeyShift) {
			u.Selected = false
		}
	}
}

func (g *clientGame) handleUnitMovement() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) && ebiten.IsFocused() {
		mx, my := ebiten.CursorPosition()
		tileX, tileY := g.screenToWorldTiles(mx, my)
		for _, u := range g.store.GetUnitsByPlayerId(g.playerId) {
			if !u.Selected {
				continue
			}
			moveStartAction := game.MoveStartAction{
				Type: game.MoveStartActionType,
				Payload: game.MoveStartPayload{
					UnitId: u.Id,
					Point:  image.Pt(tileX, tileY),
				},
			}
			g.enDispatch(moveStartAction)
		}
	}
}

func (g *clientGame) updateUnits() {
	for _, u := range g.store.GetAllUnits() {
		u.Update(g.enDispatch)
	}
}

func (g *clientGame) updateVisibility() {
	visibilityMap := g.buildVisibilityMap()
	g.applyVisibilityMap(visibilityMap)
}

func (g *clientGame) buildVisibilityMap() map[image.Point]bool {
	m := make(map[image.Point]bool, 0)

	g.initializeVisibilityMap(m)
	g.updateVisibilityFromUnits(m)

	return m
}

func (g *clientGame) initializeVisibilityMap(m map[image.Point]bool) {
	for _, t := range g.screen.tiles {
		if t.Visible || t.Unit != nil {
			m[t.Point] = false
		}
	}
}

func (g *clientGame) updateVisibilityFromUnits(m map[image.Point]bool) {
	for _, unit := range g.store.GetUnitsByPlayerId(g.playerId) {
		for _, v := range unit.ISee {
			p := unit.Position.ImagePoint().Add(v)
			m[p] = true
		}
	}
}

func (g *clientGame) applyVisibilityMap(m map[image.Point]bool) {
	for p, visible := range m {
		if t, ok := g.screen.tiles[p]; ok {
			t.Visible = visible
		}
	}
}

func (g *clientGame) screenToWorld(screenX, screenY int) (worldX, worldY int) {
	worldX = screenX + g.cameraX
	worldY = screenY + g.cameraY
	return worldX, worldY
}

func (g *clientGame) screenToWorldTiles(screenX, screenY int) (tileX, tileY int) {
	tileX = (screenX + g.cameraX) / tileSize
	tileY = (screenY + g.cameraY) / tileSize
	return tileX, tileY
}

func (g *clientGame) worldToScreen(worldX, worldY int) (screenX, screenY int) {
	screenX = worldX - g.cameraX
	screenY = worldY - g.cameraY
	return screenX, screenY
}

func (g *clientGame) queueMapLoadActions(rect image.Rectangle) {
	currRect := g.screen.rect
	if rect.Min.X < currRect.Min.X {
		action := game.NewMapLoadAction(image.Rect(rect.Min.X, rect.Min.Y, currRect.Min.X, currRect.Max.Y), g.playerId)
		g.enDispatch(action)
	}
	if rect.Min.Y < currRect.Min.Y {
		action := game.NewMapLoadAction(image.Rect(rect.Min.X, rect.Min.Y, currRect.Max.X, currRect.Min.Y), g.playerId)
		g.enDispatch(action)
	}
	if rect.Max.X > currRect.Max.X {
		action := game.NewMapLoadAction(image.Rect(currRect.Min.X, currRect.Max.Y, rect.Max.X, rect.Max.Y), g.playerId)
		g.enDispatch(action)
	}
	if rect.Max.Y > currRect.Max.Y {
		action := game.NewMapLoadAction(image.Rect(currRect.Max.Y, currRect.Min.Y, rect.Max.X, rect.Max.Y), g.playerId)
		g.enDispatch(action)
	}
}

func getRect(u *game.Unit) image.Rectangle {
	screenPosition := u.Position.Mul(tileSize).ImagePoint()
	return image.Rectangle{
		Min: screenPosition,
		Max: screenPosition.Add(u.Size),
	}
}
