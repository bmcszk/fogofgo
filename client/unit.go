package main

import (
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
}

func NewUnit(position game.PF, color color.Color, width, height int) *Unit {
	return &Unit{
		Unit: game.NewUnit(position, color, width, height),
	}
}

func (u Unit) Draw(screen *ebiten.Image, cameraX, cameraY int) {
	x := u.Position.X - float64(cameraX)
	y := u.Position.Y - float64(cameraY)

	if u.Selected {
		col := color.RGBA{0, 255, 0, 255}
		ebitenutil.DrawRect(screen, x-selectedBorder, y-selectedBorder, u.Size.X+(selectedBorder*2), u.Size.Y+(selectedBorder*2), col)
	}

	ebitenutil.DrawRect(screen, x, y, u.Size.X, u.Size.Y, u.Color)
}
