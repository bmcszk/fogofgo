package main

import (
	"image/color"
	"math"

	"github.com/bmcszk/gptrts/pkg/convert"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	UnitSpeed      = 1
	selectedBorder = 2
)

type Unit struct {
	Color                color.Color
	X, Y                 float64  // current position
	TargetX, TargetY     *float64 // target position for movement
	Selected             bool
	Width, Height        float64
	VelocityX, VelocityY float64
}

func (u *Unit) MoveTo(x, y int) {
	u.TargetX = convert.ToPointer(float64(x))
	u.TargetY = convert.ToPointer(float64(y))
}

func (u *Unit) Update() error {
	if u.TargetX == nil || u.TargetY == nil {
		return nil
	}
	// Move the unit towards the target position
	dx, dy := *u.TargetX-u.X, *u.TargetY-u.Y
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist == 0 {
		u.VelocityX, u.VelocityY = 0, 0
		u.TargetX = nil
		u.TargetY = nil
	} else {
		dx, dy = dx/dist, dy/dist
		u.VelocityX, u.VelocityY = dx*UnitSpeed, dy*UnitSpeed
	}
	u.X += u.VelocityX
	u.Y += u.VelocityY

	return nil
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
