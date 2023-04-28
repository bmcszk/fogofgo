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

type Server struct {
	game    *Game
	clients map[game.PlayerIdType]*comm.Client
}

func NewServer(g *Game) *Server {
	return &Server{
		game:    g,
		clients: make(map[game.PlayerIdType]*comm.Client, 0), // connected clients,
	}
}

func main() {
	server := NewServer(NewGame(newServerStore(), world.NewWorldService()))

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
	client := comm.NewClient(ws)

	for client.Connected {
		action, err := client.HandleInMessages()
		if err != nil {
			log.Println(err)
			continue
		}
		// register new player
		if action.GetType() == game.PlayerInitActionType {
			client.PlayerId = action.GetPayload().(game.Player).Id
			s.clients[client.PlayerId] = client
		}

		// broadcast action to others
		s.broadcastOthers(client, action)

		// synchronous dispatch func
		outActions := make([]game.Action, 0, 10)
		dispatch := func(a game.Action) {
			outActions = append(outActions, a)
		}

		// action handling
		s.game.HandleAction(action, dispatch)

		// routing output actions
		for _, a := range outActions {
			if err := s.route(client, a, dispatch); err != nil {
				log.Println(err)
			}
		}
	}
}

func (s *Server) broadcastOthers(client *comm.Client, action game.Action) {
	switch action.(type) {
	case game.PlayerInitAction, game.MapLoadAction:
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

func (s *Server) broadcastAll(action game.Action) {
	for _, c := range s.clients {
		err := c.Send(action)
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *Server) route(c *comm.Client, action game.Action, dispatch game.DispatchFunc) error {
	switch a := action.(type) {
	case game.MoveStepAction, game.MoveStopAction, game.SpawnUnitAction:
		s.game.HandleAction(a, dispatch)
		s.broadcastAll(a)
	case game.PlayerInitSuccessAction:
		if err := c.Send(action); err != nil {
			return fmt.Errorf("route %w", err)
		}
	case game.MapLoadSuccessAction:
		s.game.HandleAction(a, dispatch)
		if err := c.Send(action); err != nil {
			return fmt.Errorf("route %w", err)
		}
	default:
		s.broadcastAll(a)
	}

	return nil
}
