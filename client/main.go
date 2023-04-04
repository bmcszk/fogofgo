package main

import (
	"log"
	"net/url"

	"github.com/bmcszk/gptrts/pkg/game"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	grassImage, dirtImage *ebiten.Image
)

func init() {
	// rand.Seed(time.Now().UnixNano())

	var err error
	grassImage, _, err = ebitenutil.NewImageFromFile("grass.png")
	if err != nil {
		log.Fatalf("Failed to load grass image: %v", err)
	}

	dirtImage, _, err = ebitenutil.NewImageFromFile("dirt.png")
	if err != nil {
		log.Fatalf("Failed to load dirt image: %v", err)
	}
}

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:8000", Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer ws.Close()

	dispatchFn := dispatch(ws)
	g := NewGame(dispatchFn)

	// Read messages from the server
	go func() {
		for {
			_, bytes, err := ws.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}
			action, err := game.UnmarshalAction(bytes)
			if err != nil {
				log.Fatal(err)
			}
			if err := g.HandleAction(action); err != nil {
				log.Println(err)
			}
		}
	}()

	if err := dispatchFn(game.StartClientRequestAction{
		Type: game.StartClientRequestActionType,
		Payload: game.Player{
			Id:   game.PlayerIdType(uuid.New()),
			Name: "Kicia",
		},
	}); err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("RTS Game")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func dispatch(c *websocket.Conn) func(any) error {
	return func(action any) error {
		log.Printf("dispatch %+v", action)
		if err := c.WriteJSON(action); err != nil {
			log.Println("write:", err)
		}
		return nil
	}
}
