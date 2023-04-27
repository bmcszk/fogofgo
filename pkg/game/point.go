package game

import (
	"image"
	"math"
)

var ZeroPoint = image.Pt(0, 0)

type PF struct {
	X float64
	Y float64
}

func NewPF(x, y float64) PF {
	return PF{
		X: x, Y: y,
	}
}

func ToPF(p image.Point) PF {
	return NewPF(float64(p.X), float64(p.Y))
}

func (p PF) ImagePoint() image.Point {
	return image.Pt(int(p.X), int(p.Y))
}

func (p PF) Add(p2 PF) PF {
	return NewPF(p.X+p2.X, p.Y+p2.Y)
}

func (p PF) Mul(a float64) PF {
	return NewPF(p.X*a, p.Y*a)
}

func (p PF) Step(target PF) PF {
	s := p.Round()
	target = target.Round()
	dx, dy := target.X-s.X, target.Y-s.Y
	if dx > 0 {
		dx = 1
	} else if dx < 0 {
		dx = -1
	}
	if dy > 0 {
		dy = 1
	} else if dy < 0 {
		dy = -1
	}
	return NewPF(s.X+dx, s.Y+dy)
}

func (p PF) Round() PF {
	return NewPF(math.Round(p.X), math.Round(p.Y))
}

func (p PF) Ints() (int, int) {
	p = p.Round()
	return int(p.X), int(p.Y)
}

func (p PF) Dist(target PF) float64 {
	dx, dy := target.X-p.X, target.Y-p.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func Dist(p1 image.Point, p2 image.Point) float64 {
	dx, dy := float64(p2.X-p1.X), float64(p2.Y-p1.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

func NextStep(s image.Point, target image.Point) image.Point {
	dx, dy := target.X-s.X, target.Y-s.Y
	if dx > 0 {
		dx = 1
	} else if dx < 0 {
		dx = -1
	}
	if dy > 0 {
		dy = 1
	} else if dy < 0 {
		dy = -1
	}
	return image.Pt(s.X+dx, s.Y+dy)
}
