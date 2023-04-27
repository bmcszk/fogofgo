package game

import (
	"encoding/json"
	"errors"
	"image"

	"github.com/bmcszk/gptrts/pkg/world"
)

type ActionType string

const (
	PlayerInitActionType        ActionType = "PlayerInit"
	PlayerInitSuccessActionType ActionType = "PlayerInitSuccess"
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

type PlayerInitAction = GenericAction[Player]

type PlayerInitSuccessAction = GenericAction[PlayerInitSuccessPayload]

type PlayerInitSuccessPayload struct {
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
	var msg GenericAction[any]
	if err := json.Unmarshal(bytes, &msg); err != nil {
		return nil, err
	}
	switch msg.Type {
	case PlayerInitActionType:
		var action PlayerInitAction
		if err := json.Unmarshal(bytes, &action); err != nil {
			return nil, err
		}
		return action, nil

	case PlayerInitSuccessActionType:
		var action PlayerInitSuccessAction
		if err := json.Unmarshal(bytes, &action); err != nil {
			return nil, err
		}
		return action, nil

	case SpawnUnitActionType:
		var action SpawnUnitAction
		if err := json.Unmarshal(bytes, &action); err != nil {
			return nil, err
		}
		return action, nil

	case MoveStartActionType:
		var action MoveStartAction
		if err := json.Unmarshal(bytes, &action); err != nil {
			return nil, err
		}
		return action, nil

	case MoveStepActionType:
		var action MoveStepAction
		if err := json.Unmarshal(bytes, &action); err != nil {
			return nil, err
		}
		return action, nil

	case MoveStopActionType:
		var action MoveStopAction
		if err := json.Unmarshal(bytes, &action); err != nil {
			return nil, err
		}
		return action, nil

	case MapLoadActionType:
		var action MapLoadAction
		if err := json.Unmarshal(bytes, &action); err != nil {
			return nil, err
		}
		return action, nil

	case MapLoadSuccessActionType:
		var action MapLoadSuccessAction
		if err := json.Unmarshal(bytes, &action); err != nil {
			return nil, err
		}
		return action, nil

	default:
		return nil, errors.New("action type unrecognized")
	}
}
