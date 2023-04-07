package main

import (
	"log"
	"net/http"

	"github.com/bmcszk/gptrts/pkg/game"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Server struct {
	game      *Game
	clients   map[*websocket.Conn]bool
	broadcast chan game.Action
}

func NewServer() *Server {
	return &Server{
		clients:   make(map[*websocket.Conn]bool), // connected clients,
		broadcast: make(chan game.Action, 10),
	}
}

func main() {
	server := NewServer()

	server.game = NewGame(server.dispatch())
	// Configure websocket route
	http.HandleFunc("/ws", server.handleConnections)

	go server.handleOutMessages()

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

	s.handleInMessages(ws)
}

func (s *Server) handleInMessages(ws *websocket.Conn) {
	for {
		_, bytes, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
		}
		action, err := game.UnmarshalAction(bytes)
		if err != nil {
			log.Println(err)
		}
		log.Printf("handle %s", action.GetType())
		if err := s.game.HandleAction(action, ws); err != nil {
			log.Println(err)
		}
		if action.GetType() != game.PlayerInitActionType {
			if err := s.Broadcast(action); err != nil {
				log.Println(err)
			}
		}
	}
}

func (s *Server) handleOutMessages() {
	for action := range s.broadcast {
		if action.GetType() == game.PlayerInitSuccessActionType {
			payload := action.GetPayload().(game.PlayerInitSuccessPayload)
			if err := s.game.Players[payload.PlayerId].ws.WriteJSON(action); err != nil {
				log.Println(err)
			}
		} else {
			for ws := range s.clients {
				if err := ws.WriteJSON(action); err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func (s *Server) Broadcast(action game.Action, excludes ...*websocket.Conn) error {
	/* recipients := []*websocket.Conn{}
	for c, connected := range s.clients {
		if connected && !slices.Contains(excludes, c) {
			recipients = append(recipients, c)
		}
	}
	for _, recipient := range recipients {
		if err := recipient.WriteJSON(action); err != nil {
			return err
		}
	} */
	s.broadcast <- action
	return nil
}

func (s *Server) dispatch() game.DispatchFunc {
	return func(action game.Action) {
		log.Printf("dispatch %s", action.GetType())
		s.broadcast <- action
	}
}
