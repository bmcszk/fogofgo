package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bmcszk/gptrts/pkg/comm"
	"github.com/bmcszk/gptrts/pkg/game"
	"github.com/bmcszk/gptrts/pkg/world"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type server struct {
	game    *serverGame
	clients map[game.PlayerIdType]*comm.Client
}

func newServer(g *serverGame) *server {
	return &server{
		game:    g,
		clients: make(map[game.PlayerIdType]*comm.Client, 0), // connected clients,
	}
}

func main() {
	s := newServer(newServerGame(game.NewStoreImpl(), world.NewWorldService()))

	// Configure websocket route
	http.HandleFunc("/ws", s.handleConnections)

	// Start the server on localhost port 8000 and log any errors
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (s *server) handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	// Make sure we close the connection when the function returns
	defer func() {
		if err := ws.Close(); err != nil {
			log.Printf("Error closing websocket: %v", err)
		}
	}()

	// Register our new client
	client := comm.NewClient(ws)

	for client.Connected {
		action, err := client.HandleInMessages()
		if err != nil {
			log.Println(err)
			continue
		}
		s.processAction(client, action)
	}
}

func (s *server) processAction(client *comm.Client, action game.Action) {
	// register new player
	if action.GetType() == game.PlayerJoinActionType {
		client.PlayerId = action.GetPayload().(game.Player).Id
		s.clients[client.PlayerId] = client
	}

	// broadcast action to others
	s.broadcastOthers(client, action)

	// synchronous dispatch func
	dispatch := func(a game.Action) {
		if err := s.route(client, a); err != nil {
			log.Println(err)
		}
	}

	// action handling
	s.game.HandleAction(action, dispatch)
}

func (s *server) broadcastOthers(client *comm.Client, action game.Action) {
	switch action.(type) {
	case game.PlayerJoinAction, game.MapLoadAction:
		return
	}
	for _, c := range s.clients {
		if c != client {
			err := c.Send(action)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (s *server) broadcastAll(action game.Action) {
	for _, c := range s.clients {
		err := c.Send(action)
		if err != nil {
			log.Println(err)
		}
	}
}

// route - handler of outgoing actions
func (s *server) route(c *comm.Client, action game.Action) error {
	dispatch := func(a game.Action) {
		if err := s.route(c, a); err != nil {
			log.Println(err)
		}
	}
	switch a := action.(type) {
	case game.MoveStepAction, game.MoveStopAction, game.SpawnUnitAction:
		s.broadcastAll(a)
		s.game.HandleAction(a, dispatch)
	case game.PlayerJoinSuccessAction:
		if err := c.Send(action); err != nil {
			return fmt.Errorf("route %w", err)
		}
	case game.MapLoadSuccessAction:
		if err := c.Send(action); err != nil {
			return fmt.Errorf("route %w", err)
		}
		s.game.HandleAction(a, dispatch)
	default:
		s.broadcastAll(a)
	}

	return nil
}
