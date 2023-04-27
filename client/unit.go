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

func (u *Unit) Update() {
	u.Unit.Update()
	u.ScreenPosition = u.Position.Mul(tileSize)
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
		ebitenutil.DrawRect(screen, x-selectedBorder, y-selectedBorder, float64(u.Size.X+(selectedBorder*2)), float64(u.Size.Y+(selectedBorder*2)), col)
	}

	ebitenutil.DrawRect(screen, x, y, float64(u.Size.X), float64(u.Size.Y), u.Color)
}

func (u *Unit) GetRect() image.Rectangle {
	return image.Rectangle{
		Min: u.ScreenPosition.ImagePoint(),
		Max: u.ScreenPosition.Add(game.ToPF(u.Size)).ImagePoint(),
	}
}
