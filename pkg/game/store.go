package game

import (
	"image"

	"github.com/bmcszk/gptrts/pkg/world"
)

type Store interface {
	StoreUnit(unit Unit)
	GetUnitById(id UnitIdType) *Unit

	StorePlayer(player Player)

	GetTilesByUnitId(id UnitIdType) []*Tile

	StoreTile(tile world.Tile) *Tile

	GetTile(image.Point) *Tile
}
