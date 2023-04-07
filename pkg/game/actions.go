package game

import (
	"encoding/json"
	"errors"
)

type ActionType string

const (
	PlayerInitActionType        ActionType = "PlayerInit"
	PlayerInitSuccessActionType ActionType = "PlayerInitSuccess"
	AddUnitActionType           ActionType = "AddUnitAction"
	MoveStartActionType         ActionType = "MoveStartAction"
	MoveStepActionType          ActionType = "MoveStepAction"
	MoveStopActionType          ActionType = "MoveStopAction"
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
	Map      Map
	Units    map[UnitIdType]Unit
	Players  map[PlayerIdType]Player
}

type AddUnitAction = GenericAction[Unit]

type MoveStartAction = GenericAction[MoveStartPayload]

type MoveStartPayload struct {
	UnitId UnitIdType
	Point  PF
}

type MoveStepAction = GenericAction[MoveStepPayload]

type MoveStepPayload struct {
	UnitId   UnitIdType
	Position PF
	Path     []PF
	Step     int
}

type MoveStopAction = GenericAction[UnitIdType]

func CreateAction(actionType ActionType) (any, error) {
	switch actionType {
	case PlayerInitActionType:
		return PlayerInitAction{
			Type: actionType,
		}, nil

	case PlayerInitSuccessActionType:
		return PlayerInitSuccessAction{
			Type: actionType,
		}, nil

	case AddUnitActionType:
		return AddUnitAction{
			Type: actionType,
		}, nil

	case MoveStepActionType:
		return MoveStepAction{
			Type: actionType,
		}, nil

	case MoveStopActionType:
		return MoveStopAction{
			Type: actionType,
		}, nil

	default:
		return nil, errors.New("action type unrecognized")
	}
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

	case AddUnitActionType:
		var action AddUnitAction
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

	default:
		return nil, errors.New("action type unrecognized")
	}
}

/* func Convert(a any) Action[any] {
	switch action := a.(type) {
	case Action[~]:
	}
	return Action[any]{
		Type: a.Type,
		Payload: a.Payload,
	}
} */
