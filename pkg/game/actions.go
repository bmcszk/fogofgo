package game

import (
	"encoding/json"
	"errors"
	"image"

	"github.com/bmcszk/fogofgo/pkg/world"
)

type ActionType string

const (
	PlayerJoinActionType        ActionType = "PlayerJoin"
	PlayerJoinSuccessActionType ActionType = "PlayerJoinSuccess"
	SpawnUnitActionType         ActionType = "SpawnUnit"
	MoveStartActionType         ActionType = "MoveStart"
	MoveStepActionType          ActionType = "MoveStep"
	MoveStopActionType          ActionType = "MoveStop"
	MapLoadActionType           ActionType = "MapLoad"
	MapLoadSuccessActionType    ActionType = "MapLoadSuccess"
)

type Action interface {
	GetType() ActionType
	GetPayload() any
}

type GenericAction[T any] struct {
	Type    ActionType
	Payload T
}

func (a GenericAction[T]) GetType() ActionType {
	return a.Type
}

func (a GenericAction[T]) GetPayload() any {
	return a.Payload
}

type PlayerJoinAction = GenericAction[Player]

type PlayerJoinSuccessAction = GenericAction[PlayerJoinSuccessPayload]

type PlayerJoinSuccessPayload struct {
	PlayerId PlayerIdType
	Units    []Unit
	Players  []Player
}

type SpawnUnitAction = GenericAction[Unit]

type MoveStartAction = GenericAction[MoveStartPayload]

type MoveStartPayload struct {
	UnitId UnitIdType
	Point  image.Point
}

type MoveStepAction = GenericAction[MoveStepPayload]

type MoveStepPayload struct {
	UnitId   UnitIdType
	Position PF
	Path     []image.Point
	Step     int
}

type MoveStopAction = GenericAction[UnitIdType]

type MapLoadAction = GenericAction[MapLoadPayload]

func NewMapLoadAction(rect image.Rectangle, playerId PlayerIdType) MapLoadAction {
	return MapLoadAction{
		Type: MapLoadActionType,
		Payload: MapLoadPayload{
			WorldRequest: world.WorldRequest{
				MinX: rect.Min.X,
				MinY: rect.Min.Y,
				MaxX: rect.Max.X,
				MaxY: rect.Max.Y,
			},
			PlayerId: playerId,
		},
	}
}

type MapLoadPayload struct {
	world.WorldRequest
	PlayerId PlayerIdType
}

type MapLoadSuccessAction = GenericAction[MapLoadSuccessPayload]

type MapLoadSuccessPayload struct {
	world.WorldResponse
	PlayerId PlayerIdType
}

func UnmarshalAction(bytes []byte) (Action, error) {
	actionType, err := extractActionType(bytes)
	if err != nil {
		return nil, err
	}

	return unmarshalByType(bytes, actionType)
}

func extractActionType(bytes []byte) (ActionType, error) {
	var msg GenericAction[any]
	if err := json.Unmarshal(bytes, &msg); err != nil {
		return "", err
	}
	return msg.Type, nil
}

func unmarshalByType(bytes []byte, actionType ActionType) (Action, error) {
	switch actionType {
	case PlayerJoinActionType:
		return unmarshalPlayerJoinAction(bytes)
	case PlayerJoinSuccessActionType:
		return unmarshalPlayerJoinSuccessAction(bytes)
	case SpawnUnitActionType:
		return unmarshalSpawnUnitAction(bytes)
	case MoveStartActionType:
		return unmarshalMoveStartAction(bytes)
	case MoveStepActionType:
		return unmarshalMoveStepAction(bytes)
	case MoveStopActionType:
		return unmarshalMoveStopAction(bytes)
	case MapLoadActionType:
		return unmarshalMapLoadAction(bytes)
	case MapLoadSuccessActionType:
		return unmarshalMapLoadSuccessAction(bytes)
	default:
		return nil, errors.New("action type unrecognized")
	}
}

func unmarshalPlayerJoinAction(bytes []byte) (Action, error) {
	var action PlayerJoinAction
	if err := json.Unmarshal(bytes, &action); err != nil {
		return nil, err
	}
	return action, nil
}

func unmarshalPlayerJoinSuccessAction(bytes []byte) (Action, error) {
	var action PlayerJoinSuccessAction
	if err := json.Unmarshal(bytes, &action); err != nil {
		return nil, err
	}
	return action, nil
}

func unmarshalSpawnUnitAction(bytes []byte) (Action, error) {
	var action SpawnUnitAction
	if err := json.Unmarshal(bytes, &action); err != nil {
		return nil, err
	}
	return action, nil
}

func unmarshalMoveStartAction(bytes []byte) (Action, error) {
	var action MoveStartAction
	if err := json.Unmarshal(bytes, &action); err != nil {
		return nil, err
	}
	return action, nil
}

func unmarshalMoveStepAction(bytes []byte) (Action, error) {
	var action MoveStepAction
	if err := json.Unmarshal(bytes, &action); err != nil {
		return nil, err
	}
	return action, nil
}

func unmarshalMoveStopAction(bytes []byte) (Action, error) {
	var action MoveStopAction
	if err := json.Unmarshal(bytes, &action); err != nil {
		return nil, err
	}
	return action, nil
}

func unmarshalMapLoadAction(bytes []byte) (Action, error) {
	var action MapLoadAction
	if err := json.Unmarshal(bytes, &action); err != nil {
		return nil, err
	}
	return action, nil
}

func unmarshalMapLoadSuccessAction(bytes []byte) (Action, error) {
	var action MapLoadSuccessAction
	if err := json.Unmarshal(bytes, &action); err != nil {
		return nil, err
	}
	return action, nil
}
