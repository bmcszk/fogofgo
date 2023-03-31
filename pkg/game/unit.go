package game

import (
	"image"
	"image/color"
	"math"

	"github.com/bmcszk/gptrts/pkg/convert"
)

const (
	UnitSpeed = 1
)

type Unit struct {
	Color    color.Color
	Position PF
	Target   *PF // target position for movement
	Selected bool
	Size     PF
	Velocity PF
}

func NewUnit(position PF, color color.Color, width, height int) *Unit {
	return &Unit{
		Position: position,
		Color:    color,
		Size:     NewPF(float64(width), float64(height)),
	}
}

func (u *Unit) MoveTo(x, y int) {
	u.Target = convert.ToPointer(NewPF(float64(x), float64(y)))
}

func (u *Unit) GetRect() image.Rectangle {
	return image.Rectangle{
		Min: u.Position.ImagePoint(),
		Max: u.Position.Add(u.Size).ImagePoint(),
	}
}

func (u *Unit) Update() error {
	if u.Target == nil {
		return nil
	}
	// Move the unit towards the target position
	dx, dy := float64(u.Target.X-u.Position.X), float64(u.Target.Y-u.Position.Y)
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist < 1 {
		u.Velocity = NewPF(0, 0)
		u.Target = nil
	} else {
		dx, dy = dx/dist, dy/dist
		u.Velocity = NewPF(dx*UnitSpeed, dy*UnitSpeed)
	}
	u.Position = u.Position.Add(u.Velocity)

	return nil
}
