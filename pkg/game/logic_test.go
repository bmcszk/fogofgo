package game_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/bmcszk/gptrts/pkg/game"
	"github.com/bmcszk/gptrts/pkg/world"
	"github.com/google/uuid"
)

func TestGameLogic_HandleAction_PlayerJoinSuccessAction(t *testing.T) {
	t.Helper()
	store := game.NewStoreImpl()
	logic := game.NewGameLogic(store)

	player1 := createTestPlayer("player1")
	player2 := createTestPlayer("player2")
	unit1 := createTestUnit(player1.Id, image.Pt(0, 0))
	unit2 := createTestUnit(player2.Id, image.Pt(1, 1))

	action := game.PlayerJoinSuccessAction{
		Type: game.PlayerJoinSuccessActionType,
		Payload: game.PlayerJoinSuccessPayload{
			PlayerId: player1.Id,
			Units:    []game.Unit{*unit1, *unit2},
			Players:  []game.Player{player1, player2},
		},
	}

	var dispatchedActions []game.Action
	dispatchFunc := func(action game.Action) {
		dispatchedActions = append(dispatchedActions, action)
	}

	logic.HandleAction(action, dispatchFunc)

	// Verify units were stored
	storedUnits := store.GetAllUnits()
	if len(storedUnits) != 2 {
		t.Errorf("expected 2 units, got %d", len(storedUnits))
	}

	// Verify players were stored
	storedPlayers := store.GetAllPlayers()
	if len(storedPlayers) != 2 {
		t.Errorf("expected 2 players, got %d", len(storedPlayers))
	}

	// Verify units are placed on tiles
	tile1, exists := store.GetTile(image.Pt(0, 0))
	if !exists || tile1.Unit == nil {
		t.Error("unit1 should be placed on tile at (0,0)")
	}

	tile2, exists := store.GetTile(image.Pt(1, 1))
	if !exists || tile2.Unit == nil {
		t.Error("unit2 should be placed on tile at (1,1)")
	}
}

func TestGameLogic_HandleAction_SpawnUnitAction(t *testing.T) {
	t.Helper()
	store := game.NewStoreImpl()
	logic := game.NewGameLogic(store)

	player := createTestPlayer("testplayer")
	unit := createTestUnit(player.Id, image.Pt(2, 3))

	action := game.SpawnUnitAction{
		Type:    game.SpawnUnitActionType,
		Payload: *unit,
	}

	var dispatchedActions []game.Action
	dispatchFunc := func(action game.Action) {
		dispatchedActions = append(dispatchedActions, action)
	}

	logic.HandleAction(action, dispatchFunc)

	// Verify unit was stored
	storedUnits := store.GetAllUnits()
	if len(storedUnits) != 1 {
		t.Errorf("expected 1 unit, got %d", len(storedUnits))
	}

	// Verify unit is placed correctly
	retrievedUnit := store.GetUnitById(unit.Id)
	if retrievedUnit == nil {
		t.Fatal("unit should be retrievable by ID")
	}

	if retrievedUnit.Position != unit.Position {
		t.Errorf("expected position %v, got %v", unit.Position, retrievedUnit.Position)
	}

	// Verify unit is placed on tile
	tile, exists := store.GetTile(image.Pt(2, 3))
	if !exists || tile.Unit == nil {
		t.Error("unit should be placed on tile at (2,3)")
	}
}

func TestGameLogic_HandleAction_MoveStartAction(t *testing.T) {
	t.Helper()
	store := game.NewStoreImpl()
	logic := game.NewGameLogic(store)

	player := createTestPlayer("testplayer")
	unit := createTestUnit(player.Id, image.Pt(0, 0))
	store.StoreUnit(unit)

	targetPoint := image.Pt(5, 5)
	action := game.MoveStartAction{
		Type: game.MoveStartActionType,
		Payload: game.MoveStartPayload{
			UnitId: unit.Id,
			Point:  targetPoint,
		},
	}

	var dispatchedActions []game.Action
	dispatchFunc := func(action game.Action) {
		dispatchedActions = append(dispatchedActions, action)
	}

	logic.HandleAction(action, dispatchFunc)

	// Verify unit has path set
	retrievedUnit := store.GetUnitById(unit.Id)
	if retrievedUnit == nil {
		t.Fatal("unit should exist")
	}

	if len(retrievedUnit.Path) == 0 {
		t.Error("unit should have a path after MoveStartAction")
	}

	// Check that the path ends at the target
	lastPoint := retrievedUnit.Path[len(retrievedUnit.Path)-1]
	if lastPoint != targetPoint {
		t.Errorf("expected path to end at %v, got %v", targetPoint, lastPoint)
	}
}

func TestGameLogic_HandleAction_MoveStepAction(t *testing.T) {
	t.Helper()
	store := game.NewStoreImpl()
	logic := game.NewGameLogic(store)

	player := createTestPlayer("testplayer")
	unit := createTestUnit(player.Id, image.Pt(0, 0))

	// Place unit on initial tile
	tile := store.CreateTile(image.Pt(0, 0))
	tile.Unit = unit
	store.StoreUnit(unit)

	newPosition := game.NewPF(1.5, 1.5)
	path := []image.Point{image.Pt(1, 1), image.Pt(2, 2), image.Pt(3, 3)}

	action := game.MoveStepAction{
		Type: game.MoveStepActionType,
		Payload: game.MoveStepPayload{
			UnitId:   unit.Id,
			Position: newPosition,
			Path:     path,
			Step:     0,
		},
	}

	var dispatchedActions []game.Action
	dispatchFunc := func(action game.Action) {
		dispatchedActions = append(dispatchedActions, action)
	}

	logic.HandleAction(action, dispatchFunc)

	// Verify unit position was updated
	retrievedUnit := store.GetUnitById(unit.Id)
	if retrievedUnit == nil {
		t.Fatal("unit should exist")
	}

	if retrievedUnit.Position != newPosition {
		t.Errorf("expected position %v, got %v", newPosition, retrievedUnit.Position)
	}

	if len(retrievedUnit.Path) != len(path) {
		t.Errorf("expected path length %d, got %d", len(path), len(retrievedUnit.Path))
	}

	if retrievedUnit.Step != 0 {
		t.Errorf("expected step 0, got %d", retrievedUnit.Step)
	}

	// Verify unit is placed on new tile
	newTile, exists := store.GetTile(image.Pt(1, 1))
	if !exists || newTile.Unit == nil {
		t.Error("unit should be placed on new tile")
	}

	// Verify old tile is cleared
	oldTile, _ := store.GetTile(image.Pt(0, 0))
	if oldTile.Unit != nil {
		t.Error("old tile should be cleared")
	}
}

func TestGameLogic_HandleAction_MoveStepAction_WithCollision(t *testing.T) {
	t.Helper()
	store := game.NewStoreImpl()
	logic := game.NewGameLogic(store)

	player := createTestPlayer("testplayer")
	unit1 := createTestUnit(player.Id, image.Pt(0, 0))
	unit2 := createTestUnit(player.Id, image.Pt(1, 1))

	// Place both units
	store.StoreUnit(unit1)
	store.StoreUnit(unit2)

	// Place unit2 on the next step tile
	nextStepTile := store.CreateTile(image.Pt(1, 1))
	nextStepTile.Unit = unit2

	newPosition := game.NewPF(0.5, 0.5)
	path := []image.Point{image.Pt(1, 1)} // This will collide with unit2

	action := game.MoveStepAction{
		Type: game.MoveStepActionType,
		Payload: game.MoveStepPayload{
			UnitId:   unit1.Id,
			Position: newPosition,
			Path:     path,
			Step:     0,
		},
	}

	var dispatchedActions []game.Action
	dispatchFunc := func(action game.Action) {
		dispatchedActions = append(dispatchedActions, action)
	}

	logic.HandleAction(action, dispatchFunc)

	// Verify MoveStopAction was dispatched due to collision
	if len(dispatchedActions) != 1 {
		t.Errorf("expected 1 dispatched action, got %d", len(dispatchedActions))
	}

	if len(dispatchedActions) > 0 {
		stopAction, ok := dispatchedActions[0].(game.MoveStopAction)
		if !ok {
			t.Errorf("expected MoveStopAction, got %T", dispatchedActions[0])
		} else if stopAction.Payload != unit1.Id {
			t.Errorf("expected stop action for unit %v, got %v", unit1.Id, stopAction.Payload)
		}
	}
}

func TestGameLogic_HandleAction_MoveStopAction(t *testing.T) {
	t.Helper()
	store := game.NewStoreImpl()
	logic := game.NewGameLogic(store)

	player := createTestPlayer("testplayer")
	unit := createTestUnit(player.Id, image.Pt(0, 0))

	// Set up unit with a path and step
	unit.Path = []image.Point{image.Pt(1, 1), image.Pt(2, 2)}
	unit.Step = 1
	store.StoreUnit(unit)

	action := game.MoveStopAction{
		Type:    game.MoveStopActionType,
		Payload: unit.Id,
	}

	var dispatchedActions []game.Action
	dispatchFunc := func(action game.Action) {
		dispatchedActions = append(dispatchedActions, action)
	}

	logic.HandleAction(action, dispatchFunc)

	// Verify unit's path and step were cleared
	retrievedUnit := store.GetUnitById(unit.Id)
	if retrievedUnit == nil {
		t.Fatal("unit should exist")
	}

	if len(retrievedUnit.Path) != 0 {
		t.Errorf("expected empty path, got %v", retrievedUnit.Path)
	}

	if retrievedUnit.Step != 0 {
		t.Errorf("expected step 0, got %d", retrievedUnit.Step)
	}
}

func TestGameLogic_HandleAction_MapLoadSuccessAction(t *testing.T) {
	t.Helper()
	store := game.NewStoreImpl()
	logic := game.NewGameLogic(store)

	tiles := []world.Tile{
		{
			Point:          image.Pt(0, 0),
			Value:          "grass",
			LandType:       "plains",
			BackStyleClass: "grass",
			GroundLevel:    100,
		},
		{
			Point:          image.Pt(1, 0),
			Value:          "water",
			LandType:       "water",
			BackStyleClass: "water",
			GroundLevel:    50,
		},
	}

	action := game.MapLoadSuccessAction{
		Type: game.MapLoadSuccessActionType,
		Payload: game.MapLoadSuccessPayload{
			WorldResponse: world.WorldResponse{
				Tiles: tiles,
				MinX:  0,
				MinY:  0,
				MaxX:  1,
				MaxY:  0,
			},
			PlayerId: game.PlayerIdType(uuid.New()),
		},
	}

	var dispatchedActions []game.Action
	dispatchFunc := func(action game.Action) {
		dispatchedActions = append(dispatchedActions, action)
	}

	logic.HandleAction(action, dispatchFunc)

	// Verify tiles were stored
	tile1, exists := store.GetTile(image.Pt(0, 0))
	if !exists {
		t.Error("tile at (0,0) should exist")
	} else if tile1.Value != "grass" {
		t.Errorf("expected tile value 'grass', got '%s'", tile1.Value)
	}

	tile2, exists := store.GetTile(image.Pt(1, 0))
	if !exists {
		t.Error("tile at (1,0) should exist")
	} else if tile2.Value != "water" {
		t.Errorf("expected tile value 'water', got '%s'", tile2.Value)
	}
}

func TestGameLogic_HandleAction_UnknownAction(t *testing.T) {
	t.Helper()
	store := game.NewStoreImpl()
	logic := game.NewGameLogic(store)

	// Create a custom action that's not handled
	unknownAction := game.PlayerJoinAction{
		Type: game.PlayerJoinActionType,
		Payload: game.Player{
			Id:   game.PlayerIdType(uuid.New()),
			Name: "test",
		},
	}

	var dispatchedActions []game.Action
	dispatchFunc := func(action game.Action) {
		dispatchedActions = append(dispatchedActions, action)
	}

	// This should not panic and should handle gracefully
	logic.HandleAction(unknownAction, dispatchFunc)

	// Should not have dispatched anything or caused errors
	if len(dispatchedActions) != 0 {
		t.Errorf("expected no dispatched actions for unknown action, got %d", len(dispatchedActions))
	}
}

func TestGameLogic_PlaceUnit_Success(t *testing.T) {
	t.Helper()
	store := game.NewStoreImpl()
	logic := game.NewGameLogic(store)

	player := createTestPlayer("testplayer")
	unit := createTestUnit(player.Id, image.Pt(0, 0))

	// Use reflection to access the private placeUnit method
	// Since we can't access it directly, we'll test it through HandleAction
	action := game.SpawnUnitAction{
		Type:    game.SpawnUnitActionType,
		Payload: *unit,
	}

	var dispatchedActions []game.Action
	dispatchFunc := func(action game.Action) {
		dispatchedActions = append(dispatchedActions, action)
	}

	logic.HandleAction(action, dispatchFunc)

	// Verify unit was placed successfully
	tile, exists := store.GetTile(image.Pt(0, 0))
	if !exists {
		t.Error("tile should exist after placing unit")
	}
	if tile.Unit == nil {
		t.Error("tile should contain the unit")
	}
	if tile.Unit.Id != unit.Id {
		t.Error("tile should contain the correct unit")
	}
}

func TestGameLogic_PlaceUnit_CreatesTileWhenNotExists(t *testing.T) {
	t.Helper()
	store := game.NewStoreImpl()
	logic := game.NewGameLogic(store)

	player := createTestPlayer("testplayer")
	unit := createTestUnit(player.Id, image.Pt(10, 10))

	// Verify tile doesn't exist initially
	_, exists := store.GetTile(image.Pt(10, 10))
	if exists {
		t.Error("tile should not exist initially")
	}

	action := game.SpawnUnitAction{
		Type:    game.SpawnUnitActionType,
		Payload: *unit,
	}

	var dispatchedActions []game.Action
	dispatchFunc := func(action game.Action) {
		dispatchedActions = append(dispatchedActions, action)
	}

	logic.HandleAction(action, dispatchFunc)

	// Verify tile was created
	tile, exists := store.GetTile(image.Pt(10, 10))
	if !exists {
		t.Error("tile should be created when placing unit")
	}
	if tile.Unit == nil {
		t.Error("newly created tile should contain the unit")
	}
}

func TestGameLogic_MoveStepAction_CollisionHandling(t *testing.T) {
	t.Helper()
	store := game.NewStoreImpl()
	logic := game.NewGameLogic(store)

	player := createTestPlayer("testplayer")
	unit1 := createTestUnit(player.Id, image.Pt(0, 0))
	unit2 := createTestUnit(player.Id, image.Pt(1, 1))

	// Place unit1 initially
	store.StoreUnit(unit1)
	tile1 := store.CreateTile(image.Pt(0, 0))
	tile1.Unit = unit1

	// Place unit2 on the target position
	store.StoreUnit(unit2)
	targetTile := store.CreateTile(image.Pt(1, 1))
	targetTile.Unit = unit2

	// Try to move unit1 to where unit2 is (should cause collision)
	action := game.MoveStepAction{
		Type: game.MoveStepActionType,
		Payload: game.MoveStepPayload{
			UnitId:   unit1.Id,
			Position: game.NewPF(0.5, 0.5),
			Path:     []image.Point{image.Pt(1, 1)}, // Collision position
			Step:     0,
		},
	}

	var dispatchedActions []game.Action
	dispatchFunc := func(action game.Action) {
		dispatchedActions = append(dispatchedActions, action)
	}

	logic.HandleAction(action, dispatchFunc)

	// Should dispatch MoveStopAction due to collision
	if len(dispatchedActions) != 1 {
		t.Fatalf("expected 1 dispatched action, got %d", len(dispatchedActions))
	}

	stopAction, ok := dispatchedActions[0].(game.MoveStopAction)
	if !ok {
		t.Fatalf("expected MoveStopAction, got %T", dispatchedActions[0])
	}

	if stopAction.Payload != unit1.Id {
		t.Errorf("expected stop action for unit %v, got %v", unit1.Id, stopAction.Payload)
	}

	// unit2 should still be on its tile
	if targetTile.Unit == nil || targetTile.Unit.Id != unit2.Id {
		t.Error("unit2 should still be on its original tile")
	}
}

func TestGameLogic_MoveStepAction_SameUnitCanMoveToSamePosition(t *testing.T) {
	t.Helper()
	store := game.NewStoreImpl()
	logic := game.NewGameLogic(store)

	player := createTestPlayer("testplayer")
	unit := createTestUnit(player.Id, image.Pt(0, 0))

	// Place unit initially
	store.StoreUnit(unit)
	tile := store.CreateTile(image.Pt(0, 0))
	tile.Unit = unit

	// Move unit to the same position (should not cause collision)
	action := game.MoveStepAction{
		Type: game.MoveStepActionType,
		Payload: game.MoveStepPayload{
			UnitId:   unit.Id,
			Position: game.NewPF(0.1, 0.1),          // Slightly different position
			Path:     []image.Point{image.Pt(0, 0)}, // Same tile
			Step:     0,
		},
	}

	var dispatchedActions []game.Action
	dispatchFunc := func(action game.Action) {
		dispatchedActions = append(dispatchedActions, action)
	}

	logic.HandleAction(action, dispatchFunc)

	// Should not dispatch any stop actions
	for _, action := range dispatchedActions {
		if _, isStop := action.(game.MoveStopAction); isStop {
			t.Error("should not dispatch MoveStopAction when unit moves to same tile")
		}
	}

	// Unit should still be on the tile
	if tile.Unit == nil || tile.Unit.Id != unit.Id {
		t.Error("unit should still be on its tile")
	}
}

// Helper functions
func createTestPlayer(name string) game.Player {
	return game.Player{
		Id:    game.PlayerIdType(uuid.New()),
		Name:  name,
		Color: color.RGBA{255, 0, 0, 255},
		Start: game.NewPF(0, 0),
	}
}

func createTestUnit(owner game.PlayerIdType, position image.Point) *game.Unit {
	return &game.Unit{
		Id:       game.UnitIdType(uuid.New()),
		Owner:    owner,
		Color:    color.RGBA{0, 255, 0, 255},
		Position: game.NewPF(float64(position.X), float64(position.Y)),
		Size:     image.Pt(16, 16),
		Selected: false,
		Path:     []image.Point{},
		Step:     0,
		ISee:     []image.Point{},
	}
}
