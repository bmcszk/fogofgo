package game_test

import (
	"image"
	"testing"

	"github.com/bmcszk/gptrts/pkg/game"
)

func TestNewPF_Simple(t *testing.T) {
	p := game.NewPF(3.5, 4.2)
	if p.X != 3.5 || p.Y != 4.2 {
		t.Errorf("expected (3.5, 4.2), got (%f, %f)", p.X, p.Y)
	}
}

func TestPF_Add_Simple(t *testing.T) {
	p1 := game.NewPF(1.0, 2.0)
	p2 := game.NewPF(3.0, 4.0)
	result := p1.Add(p2)

	if result.X != 4.0 || result.Y != 6.0 {
		t.Errorf("expected (4.0, 6.0), got (%f, %f)", result.X, result.Y)
	}
}

func TestDist_Simple(t *testing.T) {
	p1 := image.Pt(0, 0)
	p2 := image.Pt(3, 4)

	distance := game.Dist(p1, p2)
	expected := 5.0

	if distance != expected {
		t.Errorf("expected %f, got %f", expected, distance)
	}
}

func TestNewStoreImpl_Simple(t *testing.T) {
	store := game.NewStoreImpl()
	if store == nil {
		t.Fatal("expected store, got nil")
	}

	// Test with empty store
	units := store.GetAllUnits()
	if len(units) != 0 {
		t.Errorf("expected 0 units in empty store, got %d", len(units))
	}

	players := store.GetAllPlayers()
	if len(players) != 0 {
		t.Errorf("expected 0 players in empty store, got %d", len(players))
	}
}

func TestActionTypes_Simple(t *testing.T) {
	if game.PlayerJoinActionType != "PlayerJoin" {
		t.Errorf("expected PlayerJoin, got %s", game.PlayerJoinActionType)
	}

	if game.SpawnUnitActionType != "SpawnUnit" {
		t.Errorf("expected SpawnUnit, got %s", game.SpawnUnitActionType)
	}
}

func TestNewPlayerId(t *testing.T) {
	id1 := game.NewPlayerId()
	id2 := game.NewPlayerId()

	if id1 == id2 {
		t.Error("expected different player IDs")
	}
}

func TestNewUnitId(t *testing.T) {
	id1 := game.NewUnitId()
	id2 := game.NewUnitId()

	if id1 == id2 {
		t.Error("expected different unit IDs")
	}
}
