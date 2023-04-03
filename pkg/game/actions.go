package game

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

type ActionType string

const (
	StartClientRequestActionType  ActionType = "StartClientRequestAction"
	StartClientResponseActionType  ActionType = "StartClientResponseAction"
	AddUnitActionType  ActionType = "AddUnitAction"
	MoveUnitActionType ActionType = "MoveUnitAction"
)

type Action struct {
	Type ActionType
}

type StartClientRequestAction struct {
	Type ActionType
}

type StartClientResponseAction struct {
	Type ActionType
	Actions []string
}

type AddUnitAction struct {
	Type ActionType
	Unit Unit
}

type MoveUnitAction struct {
	Type     ActionType
	UnitId   uuid.UUID
	Position PF
	Path     []PF
	Step     int
}

func UnmarshalAction(bytes []byte) (any, error) {
	var msg Action
	if err := json.Unmarshal(bytes, &msg);err != nil {
		return nil, err
	}
	switch msg.Type {
	case StartClientRequestActionType:
		var action StartClientRequestAction
		if err := json.Unmarshal(bytes, &action);err != nil {
			return nil, err
		}
		return action, nil

	case StartClientResponseActionType:
		var action StartClientResponseAction
		if err := json.Unmarshal(bytes, &action);err != nil {
			return nil, err
		}
		return action, nil

	case AddUnitActionType:
		var action AddUnitAction
		if err := json.Unmarshal(bytes, &action);err != nil {
			return nil, err
		}
		return action, nil
		
	case MoveUnitActionType:
		var action MoveUnitAction
		if err := json.Unmarshal(bytes, &action);err != nil {
			return nil, err
		}
		return action, nil

	default:
		return nil, errors.New("action type unrecognized")
	}
}