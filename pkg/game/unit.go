package game

import (
	"image"
	"image/color"
	"math"

	"github.com/google/uuid"
)

const (
	UnitSpeed = 0.1
)

var ZeroUnitId = UnitIdType(uuid.Nil)

const defaultSight = 5

var defaultISee []image.Point

func init() {
	defaultISee = make([]image.Point, 0, defaultSight)
	for x := -5; x <= 5; x++ {
		for y := -5; y <= 5; y++ {
			p := image.Pt(x, y)
			if Dist(p, ZeroPoint) <= defaultSight {
				defaultISee = append(defaultISee, image.Pt(x, y))
			}
		}
	}
}

type UnitIdType uuid.UUID

type Unit struct {
	Id       UnitIdType
	Owner    PlayerIdType
	Color    color.RGBA
	Position PF
	Size     image.Point
	Selected bool
	Velocity PF `json:"-"`
	Path     []image.Point
	Step     int
	dispatch DispatchFunc `json:"-"`
	ISee     []image.Point
}

func NewUnit(owner PlayerIdType, c color.RGBA, position PF, width, height int) *Unit {
	return &Unit{
		Id:       NewUnitId(),
		Owner:    owner,
		Color:    c,
		Position: position,
		Size:     image.Pt(width, height),
		ISee:     defaultISee,
	}
}

func NewUnitId() UnitIdType {
	return UnitIdType(uuid.New())
}

func (u *Unit) MoveTo(target image.Point) {
	if len(u.Path) > 0 && target == u.Path[len(u.Path)-1] {
		return
	}
	path := []image.Point{u.Position.ImagePoint()}
	path = plan(path, target)
	u.Path = path
	u.Step = 0
}

func (u *Unit) Set(unit Unit) {
	u.Step = unit.Step
	u.Position = unit.Position
	u.Path = unit.Path
}

func plan(path []image.Point, target image.Point) []image.Point {
	prevStep := path[len(path)-1]
	if prevStep == target {
		return path
	}
	nextStep := NextStep(prevStep, target)
	path = append(path, nextStep)
	return plan(path, target)
}

func (u *Unit) Update() {
	if len(u.Path) <= u.Step {
		return
	}
	// Move the unit towards the target position
	dx, dy := float64(u.Path[u.Step].X)-u.Position.X, float64(u.Path[u.Step].Y)-u.Position.Y
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist < 0.1 {
		u.Velocity = NewPF(0, 0)
		u.Position = ToPF(u.Path[u.Step])
		u.Step = u.Step + 1
		u.dispatchStep()
	} else {
		dx, dy = dx/dist, dy/dist
		u.Velocity = NewPF(dx*UnitSpeed, dy*UnitSpeed)
		u.Position = u.Position.Add(u.Velocity)
	}
}

func (u *Unit) dispatchStep() {
	moveAction := MoveStepAction{
		Type: MoveStepActionType,
		Payload: MoveStepPayload{
			UnitId:   u.Id,
			Position: u.Position,
			Path:     u.Path,
			Step:     u.Step,
		},
	}
	u.dispatch(moveAction)
}
