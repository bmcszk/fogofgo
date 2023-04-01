package main

import (
	"log"
	"net/url"

	"github.com/bmcszk/gptrts/pkg/game"
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
	g := NewGame()
	g.Init()

	u := url.URL{Scheme: "ws", Host: "localhost:8000", Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	go func() {
		for unit := range g.UnitEvents {
			if err := c.WriteJSON(unit); err != nil {
				log.Println("write:", err)
			}
		}
	}()

	// Read messages from the server
	go func() {
		for {
			var unit game.Unit
			err := c.ReadJSON(&unit)
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Println("unit:", unit.Id)
			g.Game.Units[unit.Id].Set(unit)
		}
	}()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("RTS Game")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
