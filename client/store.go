package main

import (
	"image"
	"sync"

	"github.com/bmcszk/gptrts/pkg/game"
	"github.com/bmcszk/gptrts/pkg/world"
)

type clientStore struct {
	unitMux   *sync.Mutex
	tilesMux  *sync.Mutex
	playerMux *sync.Mutex
	units     map[game.UnitIdType]*Unit
	tiles     map[image.Point]*Tile
	players   map[game.PlayerIdType]*game.Player
}

func newClientStore() *clientStore {
	return &clientStore{
		unitMux:   &sync.Mutex{},
		tilesMux:  &sync.Mutex{},
		playerMux: &sync.Mutex{},
		units:     make(map[game.UnitIdType]*Unit),
		tiles:     make(map[image.Point]*Tile),
		players:   make(map[game.PlayerIdType]*game.Player),
	}
}

func (s *clientStore) StoreUnit(unit game.Unit) {
	s.unitMux.Lock()
	defer s.unitMux.Unlock()
	if u, ok := s.units[unit.Id]; ok {
		u.Unit = &unit
		return
	}
	s.units[unit.Id] = NewUnit(&unit)
}

func (s *clientStore) GetUnitById(id game.UnitIdType) (unit *game.Unit) {
	s.unitMux.Lock()
	defer s.unitMux.Unlock()
	if u, ok := s.units[id]; ok {
		unit = u.Unit
	}
	return
}

func (s *clientStore) StorePlayer(player game.Player) {
	s.playerMux.Lock()
	defer s.playerMux.Unlock()
	s.players[player.Id] = &player
}

func (s *clientStore) GetTilesByUnitId(id game.UnitIdType) []*game.Tile {
	s.tilesMux.Lock()
	defer s.tilesMux.Unlock()
	r := make([]*game.Tile, 0)
	for _, t := range s.tiles {
		if t.UnitId == id {
			r = append(r, t.Tile)
		}
	}
	return r
}

func (s *clientStore) StoreTile(tile world.Tile) *game.Tile {
	s.tilesMux.Lock()
	defer s.tilesMux.Unlock()
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

func (s *clientStore) GetTile(point image.Point) (*game.Tile, bool) {
	s.tilesMux.Lock()
	defer s.tilesMux.Unlock()
	if t, ok := s.tiles[point]; ok {
		return t.Tile, true
	}
	return nil, false
}

func (s *clientStore) CreateTile(point image.Point) *game.Tile {
	return s.StoreTile(world.Tile{
		Point: point,
	})
}
