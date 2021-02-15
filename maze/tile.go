// Copyright ©️ 2021 oddstream.games

package maze

import (
	"fmt"
	"image/color"
	"log"
	"math/bits"
	"math/rand"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

// NORTH EAST SOUTH WEST bit patterns for presence of walls
const (
	NORTH = 0b0001 // 1 << iota
	EAST  = 0b0010 // 1 << 1
	SOUTH = 0b0100 // 1 << 2
	WEST  = 0b1000 // 1 << 3
	MASK  = 0b1111
)

var (
	reachableImages   map[uint]*ebiten.Image
	unreachableImages map[uint]*ebiten.Image
	dotImage          *ebiten.Image
	overSize          float64
	halfTileSize      float64
	wallbits          = [4]uint{NORTH, EAST, SOUTH, WEST} // map a direction (0..3) to it's bits
	wallopps          = [4]uint{SOUTH, WEST, NORTH, EAST} // map a direction (0..3) to it's opposite bits
	oppdirs           = [4]int{2, 3, 0, 1}
)

func init() {
	if 0 == TileSize {
		log.Fatal("Tile dimensions not set")
	}

	// var makeFunc func(uint) image.Image = makeTile
	reachableImages = make(map[uint]*ebiten.Image, 16)
	for i := uint(0); i < 16; i++ {
		img := makeTileImage(i, false)
		reachableImages[i] = ebiten.NewImageFromImage(img)
	}
	unreachableImages = make(map[uint]*ebiten.Image, 16)
	for i := uint(0); i < 16; i++ {
		img := makeTileImage(i, true)
		unreachableImages[i] = ebiten.NewImageFromImage(img)
	}

	// the tiles are all the same size, so pre-calc some useful variables
	actualTileSize, _ := reachableImages[0].Size()
	halfTileSize = float64(actualTileSize) / 2
	overSize = float64((actualTileSize - TileSize) / 2)

	{
		mid := float64(actualTileSize / 2)
		dc := gg.NewContext(actualTileSize, actualTileSize)
		dc.SetRGB(1, 1, 0)
		dc.DrawCircle(mid, mid, 3)
		dc.Fill()
		dc.Stroke()
		dotImage = ebiten.NewImageFromImage(dc.Image())
	}
}

// Tile object describes a tile
type Tile struct {
	// members that do not change until a new grid is created
	X, Y           int
	worldX, worldY float64 // position of tile
	edges          [4]*Tile

	// members that may change
	walls uint

	// volatile members
	visited bool
	marked  bool
	pen     bool
	parent  *Tile
}

// NewTile creates a new Tile object and returns a pointer to it
// all new tiles start with all four walls before they are carved later
func NewTile(x, y int) *Tile {
	t := &Tile{X: x, Y: y, walls: MASK}
	// worldX, worldY will be (re)set by Layout()
	return t
}

// Reset prepares a Tile for a new level by resetting just gameplay data, not structural data
func (t *Tile) Reset() {
	t.walls = MASK
	t.visited = false
	t.parent = nil
}

// Rect gives the x,y screen coords of the tile's top left and bottom right corners
func (t *Tile) Rect() (x0 int, y0 int, x1 int, y1 int) {
	x0 = t.X * TileSize
	y0 = t.Y * TileSize
	x1 = x0 + TileSize
	y1 = y0 + TileSize
	return // using named return parameters
}

// Neighbour returns the neighbouring tile in that direction
func (t *Tile) Neighbour(d int) *Tile {
	return t.edges[d]
}

// IsWall returns true if there is a wall in that direction
func (t *Tile) IsWall(d int) bool {
	bit := wallbits[d]
	return t.walls&bit == bit
}

func (t *Tile) addWall(d int) {
	t.walls |= wallbits[d]
	if tn := t.Neighbour(d); tn != nil {
		tn.walls |= wallopps[d]
	}
}

func (t *Tile) removeWall(d int) {
	if tn := t.Neighbour(d); tn != nil {
		var mask uint
		mask = MASK & (^wallbits[d])
		t.walls &= mask
		mask = MASK & (^wallopps[d])
		tn.walls &= mask
	}
}

func (t *Tile) toggleWall(d int) {
	if tn := t.Neighbour(d); tn != nil {
		if t.IsWall(d) {
			t.removeWall(d)
		} else {
			t.addWall(d)
		}
	}
}

func (t *Tile) removeAllWalls() {
	for d := 0; d < 4; d++ {
		t.removeWall(d)
	}
}

func (t *Tile) wallCount() int {
	return bits.OnesCount(t.walls)
	// count := 0
	// for i := 0; i < len(bits); i++ {
	// 	if t.walls&bits[i] == bits[i] {
	// 		count++
	// 	}
	// }
	// return count
}

func (t *Tile) fillCulDeSac() {
	if t.wallCount() == 3 {
		for d := 0; d < 4; d++ {
			if t.walls&wallbits[d] == 0 {
				t.walls = MASK
				if tn := t.Neighbour(d); tn != nil {
					tn.walls |= wallopps[d]
				}
				break
			}
		}
	}
}

func (t *Tile) recursiveBacktracker() {
	// dirs := [4]int{0, 1, 2, 3}
	// rand.Shuffle(len(dirs), func(i, j int) { dirs[i], dirs[j] = dirs[j], dirs[i] })
	dirs := rand.Perm(4)
	for d := 0; d < 4; d++ {
		dir := dirs[d]
		tn := t.Neighbour(dir)
		if tn != nil && tn.visited == false {
			t.removeWall(dir)
			tn.visited = true
			tn.recursiveBacktracker()
		}
	}
}

// Position of this tile (top left origin) in screen coords
func (t *Tile) Position() (float64, float64) {
	// Tile.Layout() may not have been called yet
	return float64(t.X * TileSize), float64(t.Y * TileSize)
	// return t.worldX, t.worldY
}

// String representation of this tile
func (t *Tile) String() string {
	return fmt.Sprintf("[%v,%v]", t.X, t.Y)
}

// AllTiles applies a func to all tiles
// func (t *Tile) AllTiles(fn func(*Tile)) {
// 	t.G.AllTiles(fn)
// }

// Layout the tile
func (t *Tile) Layout() {
	t.worldX = float64(t.X * TileSize)
	t.worldY = float64(t.Y * TileSize)
}

// Update the tile state (transitions, user input)
func (t *Tile) Update() error {
	return nil
}

func (t *Tile) debugText(screen *ebiten.Image, str string) {
	bound, _ := font.BoundString(TheAcmeFonts.large, str)
	w := (bound.Max.X - bound.Min.X).Ceil()
	h := (bound.Max.Y - bound.Min.Y).Ceil()
	x, y := t.worldX-overSize, t.worldY-overSize
	tx := int(x) + (TileSize-w)/2
	ty := int(y) + (TileSize-h)/2 + h
	c := color.RGBA{R: 0xff - colorBackground.R, G: 0xff - colorBackground.G, B: 0xff - colorBackground.B, A: 0xff}
	// var c color.Color = BasicColors["Black"]
	// ebitenutil.DrawRect(screen, float64(tx), float64(ty), float64(w), float64(h), c)
	text.Draw(screen, str, TheAcmeFonts.large, tx, ty, c)
}

// Draw renders a Tile object
func (t *Tile) Draw(screen *ebiten.Image) {

	// if float64(TileSize)+t.worldX+CameraX < 0 {
	// 	return
	// }
	// if float64(TileSize)+t.worldY+CameraY < 0 {
	// 	return
	// }

	// screenWidth, screenHeight := screen.Size()

	// if t.worldX+CameraX > float64(screenWidth) {
	// 	return
	// }
	// if t.worldY+CameraY > float64(screenHeight) {
	// 	return
	// }

	// a separate unreachable image had performance hit on big grids

	// scale, point translation, rotate, object translation

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(t.worldX-overSize, t.worldY-overSize)
	op.GeoM.Translate(CameraX, CameraY)

	// Reset RGB (not Alpha) forcibly
	// tilesheet already has black shapes
	{
		var r, g, b float64
		// op.ColorM.Scale(0, 0, 0, 1)
		r = float64(TheGrid.colorWall.R) / 0xff
		g = float64(TheGrid.colorWall.G) / 0xff
		b = float64(TheGrid.colorWall.B) / 0xff
		op.ColorM.Translate(r, g, b, 0)
	}

	if t.visited {
		screen.DrawImage(reachableImages[t.walls], op)
	} else {
		screen.DrawImage(unreachableImages[t.walls], op)

	}

	if t.marked {
		// TODO this is too slow in larger grids
		// op.ColorM.Scale(0, 0, 0, 1)
		// op.ColorM.Translate(1, 1, 0, 0)
		screen.DrawImage(dotImage, op)
	}

	// if DebugMode {
	// ebitenutil.DrawLine is really slow
	// 	if t.Y != 0 {
	// 		ebitenutil.DrawLine(screen,
	// 			CameraX+t.worldX,
	// 			CameraY+t.worldY,
	// 			CameraX+t.worldX+float64(TileSize),
	// 			CameraY+t.worldY,
	// 			BasicColors["Black"])
	// 	}
	// 	if t.X != 0 {
	// 		ebitenutil.DrawLine(screen,
	// 			CameraX+t.worldX,
	// 			CameraY+t.worldY,
	// 			CameraX+t.worldX,
	// 			CameraY+t.worldY+float64(TileSize),
	// 			BasicColors["Black"])
	// 	}
	// }
	// t.debugText(gridImage, fmt.Sprintf("%04b", t.walls))
}
