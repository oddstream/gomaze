// Copyright ©️ 2021 oddstream.games

package maze

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TilesAcross and TilesDown are package-level variables so they can be seen by Tile
var (
	TilesAcross int
	TilesDown   int
	LeftMargin  int
	TopMargin   int
	TileSize    int
)

// Grid is an object representing the grid of tiles
type Grid struct {
	tiles           []*Tile // a slice (not array!) of pointers to Tile objects
	palette         Palette
	colors          []*color.RGBA // a slice of pointers to colors for the tiles, one color per section
	colorBackground color.RGBA
	stroke          *Stroke
}

// NewGrid create a Grid object
func NewGrid(w, h int) *Grid {

	var screenWidth, screenHeight int

	if runtime.GOARCH == "wasm" {
		screenWidth, screenHeight = WindowWidth, WindowHeight
	} else {
		screenWidth, screenHeight = ebiten.WindowSize()
	}

	if w == 0 || h == 0 {
		TileSize = 100
		TilesAcross, TilesDown = screenWidth/TileSize, screenHeight/TileSize
	} else {
		possibleW := screenWidth / (w + 1) // add 1 to create margin for endcaps
		possibleW /= 20
		possibleW *= 20
		possibleH := screenHeight / (h + 1)
		possibleH /= 20
		possibleH *= 20
		// golang gotcha there isn't a vanilla math.MinInt()
		if possibleW < possibleH {
			TileSize = possibleW
		} else {
			TileSize = possibleH
		}
		TilesAcross, TilesDown = w, h
	}
	LeftMargin = (screenWidth - (TilesAcross * TileSize)) / 2
	TopMargin = (screenHeight - (TilesDown * TileSize)) / 2

	g := &Grid{tiles: make([]*Tile, TilesAcross*TilesDown)}
	for i := range g.tiles {
		g.tiles[i] = NewTile(g, i%TilesAcross, i/TilesAcross)
	}

	// link the tiles together to avoid all that tedious 2d array stuff
	for _, t := range g.tiles {
		x := t.X
		y := t.Y
		t.edges[0] = g.findTile(x, y-1)
		t.edges[1] = g.findTile(x+1, y)
		t.edges[2] = g.findTile(x, y+1)
		t.edges[3] = g.findTile(x-1, y)
	}

	InitTile()

	g.CreateNextLevel()

	return g
}

// Size returns the size of the grid in pixels
func (g *Grid) Size() (int, int) {
	return TilesAcross * TileSize, TilesDown * TileSize
}

func (g *Grid) findTile(x, y int) *Tile {

	// we re-order the tiles when dragging, to put the dragged tile at the top of the z-order
	// so we can't use i := x + (y * TilesAcross) to find index of tile in slice
	// except if this func is just used after tiles have been created

	for _, t := range g.tiles {
		if t.X == x && t.Y == y {
			return t
		}
	}
	return nil

	// if x < 0 || x >= TilesAcross {
	// 	return nil
	// }
	// if y < 0 || y >= TilesDown {
	// 	return nil
	// }
	// i := x + (y * TilesAcross)
	// if i < 0 || i > len(g.tiles) {
	// 	log.Fatal("findTile index out of bounds")
	// }
	// return g.tiles[i]
}

// findTileAt finds the tile under the mouse click or touch
func (g *Grid) findTileAt(pt image.Point) *Tile {
	for _, t := range g.tiles {
		if InRect(pt, t.Rect) {
			return t
		}
	}
	return nil
}

func (g *Grid) randomTile() *Tile {
	i := rand.Intn(len(g.tiles))
	return g.tiles[i]
}

func (g *Grid) carve() {
	t := g.randomTile()
	t.recursiveBacktracker()
}

// CreateNextLevel resets game data and moves the puzzle to the next level
func (g *Grid) CreateNextLevel() {
	for _, t := range g.tiles {
		t.Reset()
	}

	rand.Seed(262118)

	g.palette = Palettes[rand.Int()%len(Palettes)]
	g.colorBackground = CalcBackgroundColor(g.palette)

	g.carve()

	for _, t := range g.tiles {
		t.SetImage()
	}
}

// Layout implements ebiten.Game's Layout.
func (g *Grid) Layout(outsideWidth, outsideHeight int) (int, int) {
	LeftMargin = (outsideWidth - (TilesAcross * TileSize)) / 2
	TopMargin = (outsideHeight - (TilesDown * TileSize)) / 2
	for _, t := range g.tiles {
		t.Layout()
	}
	return outsideWidth, outsideHeight
}

// Update the board state (transitions, user input)
func (g *Grid) Update() error {

	if inpututil.IsKeyJustReleased(ebiten.KeyBackspace) {
		GSM.Switch(NewMenu())
	}

	for _, t := range g.tiles {
		t.Update()
	}

	return nil
}

// Draw renders the grid into the gridImage
func (g *Grid) Draw(screen *ebiten.Image) {

	screen.Fill(g.colorBackground)

	for _, t := range g.tiles {
		t.Draw(screen)
	}

	if DebugMode {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("%d,%d grid, tile size %d", TilesAcross, TilesDown, TileSize))
	}
}
