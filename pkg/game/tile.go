package game

import (
	"github.com/bmcszk/gptrts/pkg/world"
)

type Tile struct {
	*world.Tile
	Unit    *Unit
	Visible bool
}
