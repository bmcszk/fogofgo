package main

import (
	"image/color"

	"github.com/bmcszk/gptrts/pkg/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	UnitSpeed      = 1
	selectedBorder = 2
)

type Unit struct {
	game.Unit
}

func (u *Unit) Draw(screen *ebiten.Image, cameraX, cameraY int) {
	x := u.X - float64(cameraX)
	y := u.Y - float64(cameraY)

	if u.Selected {
		col := color.RGBA{0, 255, 0, 255}
		ebitenutil.DrawRect(screen, x-selectedBorder, y-selectedBorder, u.Width+(selectedBorder*2), u.Height+(selectedBorder*2), col)
	}

	ebitenutil.DrawRect(screen, x, y, u.Width, u.Height, u.Color)
}
