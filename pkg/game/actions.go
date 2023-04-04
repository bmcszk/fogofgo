package game

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

type ActionType string

const (
	StartClientRequestActionType  ActionType = "StartClientRequestAction"
	StartClientResponseActionType ActionType = "StartClientResponseAction"
	AddUnitActionType             ActionType = "AddUnitAction"
	MoveUnitActionType            ActionType = "MoveUnitAction"
	StopUnitActionType            ActionType = "StopUnitAction"
)

type Action[T any] struct {
	Type ActionType
	Payload T
}

type StartClientRequestAction struct {
	Type ActionType
}

type StartClientResponseAction struct {
	Type    ActionType
	Payload StartClientResponsePayload
}

type StartClientResponsePayload struct {
	Map      Map
	Units    map[uuid.UUID]Unit
}

type AddUnitAction struct {
	Type ActionType
	Payload Unit
}

type MoveUnitAction struct {
	Type     ActionType
	Payload MoveUnitActionPayload
}

type MoveUnitActionPayload struct {
	UnitId   uuid.UUID
	Position PF
	Path     []PF
	Step     int
}

type StopUnitAction struct {
	Type   ActionType
	UnitId uuid.UUID
}

func CreateAction(actionType ActionType) (any, error) {
	switch actionType {
	case StartClientRequestActionType:
		return StartClientRequestAction{
			Type: actionType,
		}, nil

	case StartClientResponseActionType:
		return StartClientResponseAction{
			Type: actionType,
		}, nil

	case AddUnitActionType:
		return AddUnitAction{
			Type: actionType,
		}, nil

	case MoveUnitActionType:
		return MoveUnitAction{
			Type: actionType,
		}, nil

	case StopUnitActionType:
		return StopUnitAction{
			Type: actionType,
		}, nil

	default:
		return nil, errors.New("action type unrecognized")
	}
}

func UnmarshalAction(bytes []byte) (any, error) {
	var msg Action[any]
	if err := json.Unmarshal(bytes, &msg); err != nil {
		return nil, err
	}
	switch msg.Type {
	case StartClientRequestActionType:
		var action StartClientRequestAction
		if err := json.Unmarshal(bytes, &action); err != nil {
			return nil, err
		}
		return action, nil

	case StartClientResponseActionType:
		var action StartClientResponseAction
		if err := json.Unmarshal(bytes, &action); err != nil {
			return nil, err
		}
		return action, nil

	case AddUnitActionType:
		var action AddUnitAction
		if err := json.Unmarshal(bytes, &action); err != nil {
			return nil, err
		}
		return action, nil

	case MoveUnitActionType:
		var action MoveUnitAction
		if err := json.Unmarshal(bytes, &action); err != nil {
			return nil, err
		}
		return action, nil

	case StopUnitActionType:
		var action StopUnitAction
		if err := json.Unmarshal(bytes, &action); err != nil {
			return nil, err
		}
		return action, nil

	default:
		return nil, errors.New("action type unrecognized")
	}
}
