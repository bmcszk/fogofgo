package game

import (
	"math"

	"github.com/google/uuid"
)

const (
	UnitSpeed = 0.1
)

type Unit struct {
	Id       uuid.UUID
	Position PF
	Selected bool
	Size     PF
	Velocity PF `json:"-"`
	Path     []PF
	Step     int
	Event    func(Unit) `json:"-"`
}

func NewUnit(position PF, width, height int, event func(Unit)) *Unit {
	return &Unit{
		Id:       uuid.New(),
		Position: position,
		Size:     NewPF(float64(width), float64(height)),
		Event:    event,
	}
}

func (u *Unit) MoveTo(x, y int) {
	target := NewPF(float64(x), float64(y))
	if len(u.Path) > 0 && target == u.Path[len(u.Path)-1] {
		return
	}
	path := []PF{u.Position}
	path = plan(path, target)
	u.Path = path
	u.Step = 0
}

func (u *Unit) Set(unit Unit) {
	u.Step = unit.Step
	u.Position = unit.Position
	u.Path = unit.Path
}

func plan(path []PF, target PF) []PF {
	prevStep := path[len(path)-1]
	if prevStep == target {
		return path
	}
	nextStep := prevStep.Step(target)
	path = append(path, nextStep)
	return plan(path, target)
}

func (u *Unit) Update() error {
	if len(u.Path) <= u.Step {
		return nil
	}
	// Move the unit towards the target position
	dx, dy := u.Path[u.Step].X-u.Position.X, u.Path[u.Step].Y-u.Position.Y
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist < 0.1 {
		u.Velocity = NewPF(0, 0)
		u.Position = u.Path[u.Step]
		u.Step = u.Step + 1
		u.Event(*u)
	} else {
		dx, dy = dx/dist, dy/dist
		u.Velocity = NewPF(dx*UnitSpeed, dy*UnitSpeed)
		u.Position = u.Position.Add(u.Velocity)
	}

	return nil
}
