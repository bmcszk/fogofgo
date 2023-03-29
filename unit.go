package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const UnitSpeed = 1

type Unit struct {
	Color                color.Color
	X, Y                 float64 // current position
	TargetX, TargetY     float64 // target position for movement
	Selected             bool
	Width, Height        float64
	VelocityX, VelocityY float64
}

func (u *Unit) Contains(x, y int) bool {
	return x > int(u.X-u.Width/2) &&
		x < int(u.X+u.Width/2) &&
		y > int(u.Y-u.Height/2) &&
		y < int(u.Y+u.Height/2)
}

func (u *Unit) MoveTo(x, y int) {
	u.TargetX, u.TargetY = float64(x), float64(y)
}

func (u *Unit) Update() error {
	// Move the unit towards the target position
	dx, dy := u.TargetX-u.X, u.TargetY-u.Y
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist == 0 {
		u.VelocityX, u.VelocityY = 0, 0
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
	ebitenutil.DrawRect(screen, x, y, u.Width, u.Height, u.Color)
}
