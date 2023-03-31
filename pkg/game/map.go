package game

type TileType int

const (
	Grass TileType = iota
	Dirt
)

type Tile struct {
	Type TileType
}

const (
	MapWidth  = 20
	MapHeight = 20
)

type Map struct {
	Width  int
	Height int
	Tiles  [][]Tile
}

func NewMap() *Map {
	m := &Map{Tiles: make([][]Tile, MapWidth), Width: MapWidth, Height: MapHeight}
	for x := 0; x < MapWidth; x++ {
		m.Tiles[x] = make([]Tile, MapHeight)
		for y := 0; y < MapHeight; y++ {
			if x > 10 && x < 15 && y > 10 && y < 15 {
				m.Tiles[x][y] = Tile{Type: Dirt}
			} else {
				m.Tiles[x][y] = Tile{Type: Grass}
			}
		}
	}
	return m
}
