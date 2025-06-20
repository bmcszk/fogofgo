package main

import (
	"crypto/md5"
	"image"
	"image/color"
	"log"
	"math/rand"
	"net/url"
	"os"
	"strings"

	"github.com/bmcszk/fogofgo/pkg/comm"
	"github.com/bmcszk/fogofgo/pkg/game"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	tilesImage       *ebiten.Image
	backgroundImages map[string]*ebiten.Image = make(map[string]*ebiten.Image)
)

type client struct {
	*comm.Client
	game *clientGame
}

func newClient(id game.PlayerIdType, ws *websocket.Conn) *client {
	c := comm.NewClient(ws)
	c.PlayerId = id

	return &client{
		Client: c,
	}
}

func init() {
	// Decode image from a byte slice instead of a file so that
	// this example works in any working directory.
	// If you want to use a file, there are some options:
	// 1) Use os.Open and pass the file to the image decoder.
	//    This is a very regular way, but doesn't work on browsers.
	// 2) Use ebitenutil.OpenFile and pass the file to the image decoder.
	//    This works even on browsers.
	// 3) Use ebitenutil.NewImageFromFile to create an ebiten.Image directly from a file.
	//    This also works on browsers.
	file, err := os.Open("tiles1.png")
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(file)
	// img, _, err := image.Decode(bytes.NewReader(images.Tiles_png))
	if err != nil {
		log.Fatal(err)
	}
	tilesImage = ebiten.NewImageFromImage(img)
}

func main() {
	name := getName()
	if name == "" {
		os.Exit(1)
	}

	playerId := getPlayerId(name)
	if playerId == (game.PlayerIdType{}) {
		os.Exit(1)
	}

	ws := connectToServer()
	if ws == nil {
		os.Exit(1)
	}
	defer closeWebSocket(ws)

	client := setupClient(playerId, ws)
	startMessageHandler(client)
	sendPlayerJoinAction(client, playerId, name)

	runGame(client.game, name)
}

func connectToServer() *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: "localhost:8000", Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("dial error: %v", err)
		return nil
	}
	return ws
}

func closeWebSocket(ws *websocket.Conn) {
	if err := ws.Close(); err != nil {
		log.Printf("Error closing websocket: %v", err)
	}
}

func setupClient(playerId game.PlayerIdType, ws *websocket.Conn) *client {
	c := newClient(playerId, ws)
	g := newClientGame(playerId, game.NewStoreImpl(), c.processNewAction)
	c.game = g
	return c
}

func startMessageHandler(c *client) {
	go func() {
		for c.Connected {
			action, err := c.HandleInMessages()
			if err != nil {
				log.Println(err)
				continue
			}
			c.game.HandleAction(action, c.route)
		}
	}()
}

func sendPlayerJoinAction(c *client, playerId game.PlayerIdType, name string) {
	if err := c.Send(game.PlayerJoinAction{
		Type: game.PlayerJoinActionType,
		Payload: game.Player{
			Id:    playerId,
			Name:  name,
			Color: nameToColor(name),
		},
	}); err != nil {
		log.Println(err)
	}
}

func runGame(g *clientGame, name string) {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle(name)

	if err := ebiten.RunGame(g); err != nil {
		log.Printf("Game error: %v", err)
		return
	}
}

func nameToColor(name string) color.RGBA {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	switch name {
	case "red":
		return color.RGBA{255, 0, 0, 255}
	case "green":
		return color.RGBA{0, 255, 0, 255}
	case "blue":
		return color.RGBA{0, 0, 255, 255}
	case "yellow":
		return color.RGBA{255, 255, 0, 255}
	case "cyan":
		return color.RGBA{0, 255, 255, 255}
	case "purple":
		return color.RGBA{255, 0, 255, 255}
	default:
		return randomRGBA()
	}
}

// randomRGBA generates a random color in RGBA format
func randomRGBA() color.RGBA {
	return color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255, // Set alpha to 255 for opaque color
	}
}

func getName() string {
	if len(os.Args) < 2 {
		log.Printf("Error: argument missing")
		return ""
	}
	return os.Args[1]
}

func getPlayerId(name string) game.PlayerIdType {
	hash := md5.Sum([]byte(name))
	id, err := uuid.FromBytes(hash[:])
	if err != nil {
		log.Printf("Error creating player ID: %v", err)
		return game.PlayerIdType{}
	}
	return game.PlayerIdType(id)
}

// processNewAction - handler of new actions
func (c *client) processNewAction(action game.Action) {
	if err := c.Send(action); err != nil {
		log.Println("route %w", err)
	}
	c.game.HandleAction(action, c.route)
}

// route - handler of outgoing actions
func (c *client) route(action game.Action) {
	if err := c.Send(action); err != nil {
		log.Println("route %w", err)
	}
	switch a := action.(type) {
	case game.MoveStartAction, game.MoveStepAction, game.MoveStopAction:
		c.game.HandleAction(a, c.route)
	}
}
