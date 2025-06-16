package comm

import (
	"fmt"
	"log"
	"sync"

	"github.com/bmcszk/fogofgo/pkg/game"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ws        *websocket.Conn
	Connected bool
	PlayerId  game.PlayerIdType
	mux       sync.Mutex
}

func NewClient(ws *websocket.Conn) *Client {
	c := &Client{
		ws:        ws,
		Connected: true,
		mux:       sync.Mutex{},
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
	log.Printf("player %s getting %s", uuid.UUID(c.PlayerId), action.GetType())
	return action, nil
}

func (c *Client) Send(action game.Action) error {
	if !c.Connected {
		return nil
	}
	c.mux.Lock()
	defer c.mux.Unlock()
	log.Printf("player %s sending %s", uuid.UUID(c.PlayerId), action.GetType())
	if err := c.ws.WriteJSON(action); err != nil {
		return fmt.Errorf("write %w", err)
	}
	return nil
}
