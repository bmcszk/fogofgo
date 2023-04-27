package main

import (
	"image"

	"github.com/bmcszk/gptrts/pkg/game"
	"github.com/bmcszk/gptrts/pkg/world"
)

type serverStore struct {
	units   map[game.UnitIdType]*game.Unit
	tiles   map[image.Point]*game.Tile
	players map[game.PlayerIdType]*game.Player
}

func newServerStore() *serverStore {
	return &serverStore{
		units:   make(map[game.UnitIdType]*game.Unit),
		tiles:   make(map[image.Point]*game.Tile),
		players: make(map[game.PlayerIdType]*game.Player),
	}
}

func (s *serverStore) StoreUnit(unit game.Unit) {
	s.units[unit.Id] = &unit
}

func (s *serverStore) GetUnitById(id game.UnitIdType) *game.Unit {
	return s.units[id]
}

func (s *serverStore) StorePlayer(player game.Player) {
	s.players[player.Id] = &player
}

func (s *serverStore) GetTilesByUnitId(id game.UnitIdType) []*game.Tile {
	r := make([]*game.Tile, 0)
	for _, t := range s.tiles {
		if t.UnitId == id {
			r = append(r, t)
		}
	}
	return r
}

func (s *serverStore) StoreTile(tile world.Tile) *game.Tile {
	if t, ok := s.tiles[tile.Point]; ok {
		t.Tile = &tile
	} else {
		s.tiles[tile.Point] = &game.Tile{
			Tile: &tile,
		}
	}
	return s.tiles[tile.Point]
}

func (s *serverStore) GetTile(point image.Point) *game.Tile {
	if t, ok := s.tiles[point]; ok {
		return t
	}
	return s.StoreTile(world.Tile{
		Point: point,
	})
}
