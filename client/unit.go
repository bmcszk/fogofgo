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
}

func NewUnit(u *game.Unit) *Unit {
	return &Unit{
		Unit:           u,
		ScreenPosition: u.Position.Mul(tileSize),
	}
}

func (u *Unit) Update() error {
	if err := u.Unit.Update(); err != nil {
		return err
	}
	u.ScreenPosition = u.Position.Mul(tileSize)
	return nil
}

func (u *Unit) Draw(screen *ebiten.Image, cameraX, cameraY int) {
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
