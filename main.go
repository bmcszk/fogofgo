package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480
	tileSize     = 32

	MapWidth  = 100
	MapHeight = 60
)

var (
	grassImage, dirtImage *ebiten.Image
)

type TileType int

const (
	Grass TileType = iota
	Dirt
)

type Tile struct {
	Type TileType
}

type Map struct {
	Width  int
	Height int
	Tiles  [][]Tile
}

type Unit struct {
	X, Y  int
	Color color.Color
}

type Game struct {
	Map    Map
	Units  []*Unit
	Player *Unit
}

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
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("RTS Game")

	g := &Game{}
	g.Init()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	// Move the player unit
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.Player.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.Player.Y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.Player.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.Player.X += 1
	}

	// Clamp the player unit to the map bounds
	if g.Player.X < 0 {
		g.Player.X = 0
	}
	if g.Player.X >= len(g.Map.Tiles[0])*tileSize {
		g.Player.X = len(g.Map.Tiles[0])*tileSize - 1
	}
	if g.Player.Y < 0 {
		g.Player.Y = 0
	}
	if g.Player.Y >= len(g.Map.Tiles)*tileSize {
		g.Player.Y = len(g.Map.Tiles)*tileSize - 1
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the map
	for y := range g.Map.Tiles {
		for x, tile := range g.Map.Tiles[y] {
			var img *ebiten.Image
			switch tile.Type {
			case Grass:
				img = grassImage
			case Dirt:
				img = dirtImage
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*tileSize), float64(y*tileSize))
			screen.DrawImage(img, op)
		}
	}

	// Draw the player unit
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.Player.X), float64(g.Player.Y))
	op.ColorM.Scale(1, 0.5, 0.5, 1) // Tint the unit red
	unitSize := float64(tileSize) / 2
	op.GeoM.Translate(-unitSize, -unitSize)
	op.GeoM.Scale(unitSize/16, unitSize/16)
	ebitenutil.DrawRect(screen, float64(g.Player.X), float64(g.Player.Y), unitSize, unitSize, g.Player.Color)

	// Draw the other units
	for _, unit := range g.Units {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(unit.X), float64(unit.Y))
		unitSize := float64(tileSize) / 2
		op.GeoM.Translate(-unitSize, -unitSize)
		op.GeoM.Scale(unitSize/16, unitSize/16)
		ebitenutil.DrawRect(screen, float64(unit.X), float64(unit.Y), unitSize, unitSize, unit.Color)
	}
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Calculate the desired screen size based on the size of the map
	screenWidth = len(g.Map.Tiles[0]) * tileSize
	screenHeight = len(g.Map.Tiles) * tileSize

	// Scale the screen if it is too large to fit
	if screenWidth > outsideWidth || screenHeight > outsideHeight {
		scale := math.Min(float64(outsideWidth)/float64(screenWidth), float64(outsideHeight)/float64(screenHeight))
		screenWidth = int(float64(screenWidth) * scale)
		screenHeight = int(float64(screenHeight) * scale)
	}

	return screenWidth, screenHeight
}

func (g *Game) Init() {
	// Initialize the map tiles
	g.Map = Map{Tiles: make([][]Tile, MapWidth), Width: MapWidth, Height: MapHeight}
	for x := 0; x < MapWidth; x++ {
		g.Map.Tiles[x] = make([]Tile, MapHeight)
		for y := 0; y < MapHeight; y++ {
			g.Map.Tiles[x][y] = Tile{Type: Grass}
		}
	}
	playerUnit := &Unit{
		X: 3, Y: 3, Color: color.RGBA{255, 0, 0, 255},
	}

	g.Player = playerUnit
	g.Units = append(g.Units, playerUnit)
}
