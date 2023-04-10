package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bmcszk/gptrts/pkg/comm"
	"github.com/bmcszk/gptrts/pkg/game"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Server struct {
	game    *Game
	clients []*comm.Client
}

func NewServer() *Server {
	return &Server{
		clients: make([]*comm.Client, 0), // connected clients,
	}
}

func main() {
	server := NewServer()

	g := NewGame(server.dispatch)

	server.game = g

	// Configure websocket route
	http.HandleFunc("/ws", server.handleConnections)

	// Start the server on localhost port 8000 and log any errors
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client

	actions := make(chan game.Action, 10)
	dispatch := func(a game.Action) { actions <- a }

	client := comm.NewClient(ws, dispatch)
	s.clients = append(s.clients, client)

	go func(acts <-chan game.Action) {
		for a := range acts {
			if err := s.route(client, a); err != nil {
				log.Println(err)
			}
		}
	}(actions)

	for client.Connected {
		action, err := client.HandleInMessages()
		if err != nil {
			log.Println(err)
			continue
		}
		if action.GetType() == game.PlayerInitActionType {
			client.PlayerId = action.GetPayload().(game.Player).Id
		}
		s.broadcast(client, action, dispatch)
		s.game.HandleAction(action)
	}
}

func (s *Server) broadcast(client *comm.Client, action game.Action, d func(game.Action)) {
	for _, c := range s.clients {
		if c != client {
			c.Dispatch(action)
		}
	}
}

func (s *Server) dispatch(action game.Action) {
	for _, c := range s.clients {
		c.Dispatch(action)
	}
}

func (s *Server) route(c *comm.Client, action game.Action) error {
	switch a := action.(type) {
	case game.MoveStartAction, game.MoveStepAction, game.MoveStopAction:
		s.game.HandleAction(a)
	case game.PlayerInitSuccessAction:
		if a.Payload.PlayerId != c.PlayerId {
			return nil
		}
	}

	if err := c.Send(action); err != nil {
		return fmt.Errorf("route %w", err)
	}
	return nil
}
