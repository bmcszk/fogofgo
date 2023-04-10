package main

import (
	"image"
	"image/color"

	"github.com/bmcszk/gptrts/pkg/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	selectedBorder = 2
)

type Unit struct {
	*game.Unit
	ScreenPosition game.PF
	visible        bool
}

func NewUnit(u *game.Unit) *Unit {
	return &Unit{
		Unit:           u,
		ScreenPosition: u.Position.Mul(tileSize),
	}
}

func (u *Unit) Update(playerUnits []*Unit) {
	u.Unit.Update()
	u.ScreenPosition = u.Position.Mul(tileSize)
	u.visible = u.isVisibile(playerUnits)
}

func (u *Unit) Draw(screen *ebiten.Image, cameraX, cameraY int) {
	if !u.visible {
		return
	}
	u.ScreenPosition = u.Position.Mul(tileSize)
	x := u.ScreenPosition.X - float64(cameraX)
	y := u.ScreenPosition.Y - float64(cameraY)

	if u.Selected {
		col := color.RGBA{0, 255, 0, 255}
		ebitenutil.DrawRect(screen, x-selectedBorder, y-selectedBorder, u.Size.X+(selectedBorder*2), u.Size.Y+(selectedBorder*2), col)
	}

	ebitenutil.DrawRect(screen, x, y, u.Size.X, u.Size.Y, u.Color)
}

func (u *Unit) GetRect() image.Rectangle {
	return image.Rectangle{
		Min: u.ScreenPosition.ImagePoint(),
		Max: u.ScreenPosition.Add(u.Size).ImagePoint(),
	}
}

func (u *Unit) isVisibile(playerUnits []*Unit) bool {
	for _, pu := range playerUnits {
		if u.Position.Dist(pu.Position) <= 5 {
			return true
		}
	}
	return false
}

