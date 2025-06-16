package game

import (
	"github.com/bmcszk/fogofgo/pkg/world"
)

type Tile struct {
	*world.Tile
	Unit    *Unit
	Visible bool
}
