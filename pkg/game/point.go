package game

import "image"

type PF struct {
	X float64
	Y float64
}

func NewPF(x, y float64) PF {
	return PF{
		X: x, Y: y,
	}
}

func (p PF) ImagePoint() image.Point {
	return image.Pt(int(p.X), int(p.Y))
}

func (p PF) Add(p2 PF) PF {
	return NewPF(p.X+p2.X, p.Y+p2.Y)
}
