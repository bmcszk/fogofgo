package game

import (
	"image"
	"sync"

	"github.com/bmcszk/fogofgo/pkg/world"
)

type Store interface {
	StoreUnit(unit *Unit)
	GetUnitById(id UnitIdType) *Unit
	GetAllUnits() []*Unit
	GetUnitsByPlayerId(id PlayerIdType) []*Unit

	GetPlayer(id PlayerIdType) (*Player, bool)
	GetAllPlayers() []*Player
	StorePlayer(player Player)

	GetTilesByUnitId(id UnitIdType) []*Tile
	StoreTile(tile world.Tile) *Tile
	GetTile(image.Point) (*Tile, bool)
	CreateTile(image.Point) *Tile
	GetTilesByRect(rect image.Rectangle) map[image.Point]*Tile
}

type StoreImpl struct {
	unitMux   *sync.Mutex
	tilesMux  *sync.Mutex
	playerMux *sync.Mutex
	units     map[UnitIdType]*Unit
	tiles     map[image.Point]*Tile
	players   map[PlayerIdType]*Player
}

func NewStoreImpl() *StoreImpl {
	return &StoreImpl{
		unitMux:   &sync.Mutex{},
		tilesMux:  &sync.Mutex{},
		playerMux: &sync.Mutex{},
		units:     make(map[UnitIdType]*Unit),
		tiles:     make(map[image.Point]*Tile),
		players:   make(map[PlayerIdType]*Player),
	}
}

func (s *StoreImpl) GetAllUnits() []*Unit {
	s.unitMux.Lock()
	defer s.unitMux.Unlock()
	r := make([]*Unit, 0, len(s.units))
	for _, u := range s.units {
		r = append(r, u)
	}
	return r
}

func (s *StoreImpl) GetUnitsByPlayerId(id PlayerIdType) []*Unit {
	s.unitMux.Lock()
	defer s.unitMux.Unlock()
	r := make([]*Unit, 0, len(s.units))
	for _, u := range s.units {
		if u.Owner == id {
			r = append(r, u)
		}
	}
	return r
}

func (s *StoreImpl) StoreUnit(unit *Unit) {
	s.unitMux.Lock()
	defer s.unitMux.Unlock()
	s.units[unit.Id] = unit
}

func (s *StoreImpl) GetUnitById(id UnitIdType) *Unit {
	s.unitMux.Lock()
	defer s.unitMux.Unlock()
	return s.units[id]
}

func (s *StoreImpl) GetPlayer(id PlayerIdType) (*Player, bool) {
	s.playerMux.Lock()
	defer s.playerMux.Unlock()
	if p, ok := s.players[id]; ok {
		return p, true
	}
	return nil, false
}

func (s *StoreImpl) GetAllPlayers() []*Player {
	s.playerMux.Lock()
	defer s.playerMux.Unlock()
	r := make([]*Player, 0, len(s.players))
	for _, p := range s.players {
		r = append(r, p)
	}
	return r
}

func (s *StoreImpl) StorePlayer(player Player) {
	s.playerMux.Lock()
	defer s.playerMux.Unlock()
	s.players[player.Id] = &player
}

func (s *StoreImpl) GetTilesByUnitId(id UnitIdType) []*Tile {
	s.tilesMux.Lock()
	defer s.tilesMux.Unlock()
	r := make([]*Tile, 0)
	for _, t := range s.tiles {
		if t.Unit != nil && t.Unit.Id == id {
			r = append(r, t)
		}
	}
	return r
}

func (s *StoreImpl) StoreTile(tile world.Tile) *Tile {
	s.tilesMux.Lock()
	defer s.tilesMux.Unlock()
	return s.storeTileInternal(tile)
}

func (s *StoreImpl) storeTileInternal(tile world.Tile) *Tile {
	if t, ok := s.tiles[tile.Point]; ok {
		t.Tile = &tile
	} else {
		s.tiles[tile.Point] = &Tile{
			Tile: &tile,
		}
	}
	return s.tiles[tile.Point]
}

func (s *StoreImpl) GetTile(point image.Point) (*Tile, bool) {
	s.tilesMux.Lock()
	defer s.tilesMux.Unlock()
	if t, ok := s.tiles[point]; ok {
		return t, true
	}
	return nil, false
}

func (s *StoreImpl) CreateTile(point image.Point) *Tile {
	return s.StoreTile(world.Tile{
		Point: point,
	})
}

func (s *StoreImpl) GetTilesByRect(rect image.Rectangle) map[image.Point]*Tile {
	s.tilesMux.Lock()
	defer s.tilesMux.Unlock()
	size := rect.Size()
	area := size.X * size.Y
	r := make(map[image.Point]*Tile, area)
	for x := rect.Min.X; x <= rect.Max.X; x++ {
		for y := rect.Min.Y; y <= rect.Max.Y; y++ {
			p := image.Pt(x, y)
			t, ok := s.tiles[p]
			if !ok {
				t = s.storeTileInternal(world.Tile{
					Point: p,
				})
			}
			r[p] = t
		}
	}
	return r
}
