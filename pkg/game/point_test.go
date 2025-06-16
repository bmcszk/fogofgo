package game_test

import (
	"image"
	"math"
	"testing"

	"github.com/bmcszk/fogofgo/pkg/game"
)

func TestZeroPoint(t *testing.T) {
	if game.ZeroPoint.X != 0 || game.ZeroPoint.Y != 0 {
		t.Errorf("expected (0,0), got (%d,%d)", game.ZeroPoint.X, game.ZeroPoint.Y)
	}
}

func TestNewPF(t *testing.T) {
	p := game.NewPF(3.5, 4.2)
	if p.X != 3.5 || p.Y != 4.2 {
		t.Errorf("expected (3.5, 4.2), got (%f, %f)", p.X, p.Y)
	}
}

func TestToPF(t *testing.T) {
	point := image.Pt(5, 7)
	pf := game.ToPF(point)
	if pf.X != 5.0 || pf.Y != 7.0 {
		t.Errorf("expected (5.0, 7.0), got (%f, %f)", pf.X, pf.Y)
	}
}

func TestPF_ImagePoint(t *testing.T) {
	tests := []struct {
		name     string
		pf       game.PF
		expected image.Point
	}{
		{"positive values", game.NewPF(3.7, 4.2), image.Pt(3, 4)},
		{"negative values", game.NewPF(-2.8, -1.3), image.Pt(-2, -1)},
		{"zero values", game.NewPF(0.0, 0.0), image.Pt(0, 0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pf.ImagePoint()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestPF_Add(t *testing.T) {
	p1 := game.NewPF(1.5, 2.5)
	p2 := game.NewPF(2.0, 3.0)
	result := p1.Add(p2)

	expected := game.NewPF(3.5, 5.5)
	if result.X != expected.X || result.Y != expected.Y {
		t.Errorf("expected (%f, %f), got (%f, %f)", expected.X, expected.Y, result.X, result.Y)
	}
}

func TestPF_Mul(t *testing.T) {
	p := game.NewPF(2.0, 3.0)
	result := p.Mul(2.5)

	expected := game.NewPF(5.0, 7.5)
	if result.X != expected.X || result.Y != expected.Y {
		t.Errorf("expected (%f, %f), got (%f, %f)", expected.X, expected.Y, result.X, result.Y)
	}
}

func TestPF_Round(t *testing.T) {
	tests := []struct {
		name     string
		pf       game.PF
		expected game.PF
	}{
		{"round up", game.NewPF(2.7, 3.6), game.NewPF(3.0, 4.0)},
		{"round down", game.NewPF(2.3, 3.2), game.NewPF(2.0, 3.0)},
		{"round half", game.NewPF(2.5, 3.5), game.NewPF(3.0, 4.0)},
		{"already rounded", game.NewPF(2.0, 3.0), game.NewPF(2.0, 3.0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pf.Round()
			if result.X != tt.expected.X || result.Y != tt.expected.Y {
				t.Errorf("expected (%f, %f), got (%f, %f)", tt.expected.X, tt.expected.Y, result.X, result.Y)
			}
		})
	}
}

func TestPF_Ints(t *testing.T) {
	pf := game.NewPF(3.7, 4.2)
	x, y := pf.Ints()

	if x != 4 || y != 4 { // Should round to 4, 4
		t.Errorf("expected (4, 4), got (%d, %d)", x, y)
	}
}

func TestPF_Dist(t *testing.T) {
	p1 := game.NewPF(0.0, 0.0)
	p2 := game.NewPF(3.0, 4.0)

	distance := p1.Dist(p2)
	expected := 5.0 // 3-4-5 triangle

	if math.Abs(distance-expected) > 0.001 {
		t.Errorf("expected %f, got %f", expected, distance)
	}
}

func TestPF_Step(t *testing.T) {
	tests := []struct {
		name     string
		start    game.PF
		target   game.PF
		expected game.PF
	}{
		{"move right", game.NewPF(0.0, 0.0), game.NewPF(5.0, 0.0), game.NewPF(1.0, 0.0)},
		{"move left", game.NewPF(5.0, 0.0), game.NewPF(0.0, 0.0), game.NewPF(4.0, 0.0)},
		{"move up", game.NewPF(0.0, 0.0), game.NewPF(0.0, 5.0), game.NewPF(0.0, 1.0)},
		{"move down", game.NewPF(0.0, 5.0), game.NewPF(0.0, 0.0), game.NewPF(0.0, 4.0)},
		{"move diagonal", game.NewPF(0.0, 0.0), game.NewPF(3.0, 3.0), game.NewPF(1.0, 1.0)},
		{"same position", game.NewPF(2.0, 2.0), game.NewPF(2.0, 2.0), game.NewPF(2.0, 2.0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.start.Step(tt.target)
			if result.X != tt.expected.X || result.Y != tt.expected.Y {
				t.Errorf("expected (%f, %f), got (%f, %f)", tt.expected.X, tt.expected.Y, result.X, result.Y)
			}
		})
	}
}

func TestDist(t *testing.T) {
	p1 := image.Pt(0, 0)
	p2 := image.Pt(3, 4)

	distance := game.Dist(p1, p2)
	expected := 5.0 // 3-4-5 triangle

	if math.Abs(distance-expected) > 0.001 {
		t.Errorf("expected %f, got %f", expected, distance)
	}
}

func TestNextStep(t *testing.T) {
	tests := []struct {
		name     string
		start    image.Point
		target   image.Point
		expected image.Point
	}{
		{"move right", image.Pt(0, 0), image.Pt(5, 0), image.Pt(1, 0)},
		{"move left", image.Pt(5, 0), image.Pt(0, 0), image.Pt(4, 0)},
		{"move up", image.Pt(0, 0), image.Pt(0, 5), image.Pt(0, 1)},
		{"move down", image.Pt(0, 5), image.Pt(0, 0), image.Pt(0, 4)},
		{"move diagonal", image.Pt(0, 0), image.Pt(3, 3), image.Pt(1, 1)},
		{"same position", image.Pt(2, 2), image.Pt(2, 2), image.Pt(2, 2)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := game.NextStep(tt.start, tt.target)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
