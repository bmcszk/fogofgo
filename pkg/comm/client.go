package comm

import (
	"fmt"
	"log"

	"github.com/bmcszk/gptrts/pkg/game"
	"github.com/gorilla/websocket"
)

type Client struct {
	ws         *websocket.Conn
	PlayerId   game.PlayerIdType
	OutActions chan game.Action
	Connected  bool
}

func NewClient(ws *websocket.Conn) *Client {
	c := &Client{
		ws:         ws,
		OutActions: make(chan game.Action, 10),
		Connected:  true,
	}
	return c
}

func (c *Client) HandleInMessages() (game.Action, error) {
	msgType, bytes, err := c.ws.ReadMessage()
	if msgType == websocket.CloseMessage || (err != nil && err == websocket.ErrCloseSent) {
		if err := c.ws.Close(); err != nil {
			log.Println(err)
		}
		c.Connected = false
		log.Printf("player %s connection closed", c.PlayerId)
		return game.GenericAction[any]{}, err
	} else if err != nil {
		log.Println(err)
		return game.GenericAction[any]{}, err
	}
	action, err := game.UnmarshalAction(bytes)
	if err != nil {
		log.Println(err)
		return game.GenericAction[any]{}, err
	}
	log.Printf("player %s getting %s", c.PlayerId, action.GetType())
	return action, nil
}

func (c *Client) Dispatch(action game.Action) {
	log.Printf("player %s dispatch %s", c.PlayerId, action.GetType())
	c.OutActions <- action
}

func (c *Client) Send(action game.Action) error {
	log.Printf("player %s send %s", c.PlayerId, action.GetType())
	if err := c.ws.WriteJSON(action); err != nil {
		return fmt.Errorf("write %w", err)
	}
	return nil
}
