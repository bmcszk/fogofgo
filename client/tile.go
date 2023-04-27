package main

import (
	"github.com/bmcszk/gptrts/pkg/game"
)

const (
	tileSize = 16
	tileXNum = 7
)

type Tile struct {
	*game.Tile
	visible bool
}
