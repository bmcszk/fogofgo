package main

import (
	"log"
	"net/http"

	"github.com/bmcszk/gptrts/pkg/game"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Server struct {
	clients map[*websocket.Conn]bool
	game *game.Game
}

func NewServer(g *game.Game) *Server {
	return &Server{
		clients: make(map[*websocket.Conn]bool), // connected clients,
		game: g,
	}
}

func main() {
	server := NewServer(game.NewGame())
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
		var unit game.Unit
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&unit)
		if err != nil {
			log.Printf("error: %v", err)
			delete(s.clients, ws)
			break
		}
		s.handleUnit(unit)
	}
}

func (s *Server) handleUnit(unit game.Unit) {
	log.Printf("Unit: %+v\n", unit)
	s.game.Units[unit.Id] = &unit
}

