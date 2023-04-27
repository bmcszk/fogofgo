package main

import (
	"image"

	"github.com/bmcszk/gptrts/pkg/game"
	"github.com/bmcszk/gptrts/pkg/world"
)

type clientStore struct {
	units   map[game.UnitIdType]*Unit
	tiles   map[image.Point]*Tile
	players map[game.PlayerIdType]*game.Player
}

func newClientStore() *clientStore {
	return &clientStore{
		units:   make(map[game.UnitIdType]*Unit),
		tiles:   make(map[image.Point]*Tile),
		players: make(map[game.PlayerIdType]*game.Player),
	}
}

func (s *clientStore) StoreUnit(unit game.Unit) {
	if u, ok := s.units[unit.Id]; ok {
		u.Unit = &unit
		return
	}
	s.units[unit.Id] = NewUnit(&unit)
}

func (s *clientStore) GetUnitById(id game.UnitIdType) (unit *game.Unit) {
	if u, ok := s.units[id]; ok {
		unit = u.Unit
	}
	return
}

func (s *clientStore) StorePlayer(player game.Player) {
	s.players[player.Id] = &player
}

func (s *clientStore) GetTilesByUnitId(id game.UnitIdType) []*game.Tile {
	r := make([]*game.Tile, 0)
	for _, t := range s.tiles {
		if t.UnitId == id {
			r = append(r, t.Tile)
		}
	}
	return r
}

func (s *clientStore) StoreTile(tile world.Tile) *game.Tile {
	if t, ok := s.tiles[tile.Point]; ok {
		if t.Tile == nil {
			t.Tile = &game.Tile{}
		}
		t.Tile.Tile = &tile
	} else {
		s.tiles[tile.Point] = &Tile{
			Tile: &game.Tile{
				Tile: &tile,
			},
		}
	}
	return s.tiles[tile.Point].Tile
}

func (s *clientStore) GetTile(point image.Point) (tile *game.Tile) {
	if t, ok := s.tiles[point]; ok {
		return t.Tile
	}
	return s.StoreTile(world.Tile{
		Point: point,
	})
}
