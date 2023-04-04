package main

import (
	"image/color"
	"log"
	"net/http"

	"github.com/bmcszk/gptrts/pkg/game"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/slices"
)

var upgrader = websocket.Upgrader{}

type Server struct {
	game     *game.Game
	clients  map[*websocket.Conn]bool
	dispatch func(any) error
}

func NewServer(g *game.Game, clients map[*websocket.Conn]bool, dispatchFn func(any) error) *Server {
	return &Server{
		game:     g,
		clients:  clients, // connected clients,
		dispatch: dispatchFn,
	}
}

func main() {
	clients := make(map[*websocket.Conn]bool)
	dispatchFn := dispatch(clients)
	g := game.NewGame(dispatchFn)

	g.HandleAction(game.AddUnitAction{
		Type: game.AddUnitActionType,
		Payload: *game.NewUnit(color.RGBA{255, 0, 0, 255}, game.NewPF(0, 0), 32, 32),
	})

	g.HandleAction(game.AddUnitAction{
		Type: game.AddUnitActionType,
		Payload: *game.NewUnit(color.RGBA{0, 0, 255, 255}, game.NewPF(1, 0), 32, 32),
	})

	server := NewServer(g, clients, dispatchFn)
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
	s.clients[ws] = true

	for {
		_, bytes, err := ws.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		action, err := game.UnmarshalAction(bytes)
		if err != nil {
			log.Fatal(err)
		}
		if err := s.game.HandleAction(action); err != nil {
			log.Println(err)
		}
		if err := s.Broadcast(action); err != nil {
			log.Println(err)
		}
	}
}

func (s *Server) Broadcast(action any, excludes ...*websocket.Conn) error {
	recipients := []*websocket.Conn{}
	for c, connected := range s.clients {
		if connected && !slices.Contains(excludes, c) {
			recipients = append(recipients, c)
		}
	}
	for _, recipient := range recipients {
		if err := recipient.WriteJSON(action); err != nil {
			return err
		}
	}
	return nil
}

func dispatch(clients map[*websocket.Conn]bool) func(any) error {
	return func(action any) error {
		log.Printf("dispatch %+v", action)
		for c, connected := range clients {
			if connected {
				if err := c.WriteJSON(action); err != nil {
					log.Println("write:", err)
				}
			}
		}
		return nil
	}
}
