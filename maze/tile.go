// Copyright ©️ 2021 oddstream.games

package maze

import (
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	tileImageLibrary map[uint]*ebiten.Image
	overSize         float64
	halfTileSize     float64
	bits             = [4]uint{NORTH, EAST, SOUTH, WEST} // map a direction (0..3) to it's bits
	opps             = [4]uint{SOUTH, WEST, NORTH, EAST} // map a direction (0..3) to it's opposite bits
	oppdirs          = [4]int{2, 3, 0, 1}
)

// InitTile used to be init(), but TileSize may not be set yet, hence this func called from NewGrid()
func InitTile() {
	if 0 == TileSize {
		log.Fatal("Tile dimensions not set")
	}

	var makeFunc func(uint, int) image.Image = makeTile
	tileImageLibrary = make(map[uint]*ebiten.Image, 16)
	for i := uint(0); i < 16; i++ {
		img := makeFunc(i, TileSize)
		tileImageLibrary[i] = ebiten.NewImageFromImage(img)
	}

	// the tiles are all the same size, so pre-calc some useful variables
	actualTileSize, _ := tileImageLibrary[0].Size()
	halfTileSize = float64(actualTileSize) / 2
	overSize = float64((actualTileSize - TileSize) / 2)
}

// Tile object describes a tile
type Tile struct {
	// members that do not change until a new grid is created
	G            *Grid
	X, Y         int
	homeX, homeY float64 // position of tile
	edges        [4]*Tile

	// members that may change
	tileImage *ebiten.Image
	walls     uint
	visited   bool
}

// NewTile creates a new Tile object and returns a pointer to it
// all new tiles start with all four walls before they are carved later
func NewTile(g *Grid, x, y int) *Tile {
	t := &Tile{G: g, X: x, Y: y, walls: MASK}
	// homeX, homeY, offsetX, offsetY will be (re)set by Layout()
	return t
}

// SetImage is used when all walls are carved
func (t *Tile) SetImage() {
	t.tileImage = tileImageLibrary[t.walls]
	if t.tileImage == nil {
		log.Fatal("tileImage is nil when walls == ", t.walls)
	}
}

// Reset prepares a Tile for a new level by resetting just gameplay data, not structural data
func (t *Tile) Reset() {
	t.tileImage = nil
	t.walls = MASK
	t.visited = false
}

// Rect gives the x,y screen coords of the tile's top left and bottom right corners
func (t *Tile) Rect() (x0 int, y0 int, x1 int, y1 int) {
	x0 = t.X*TileSize + LeftMargin
	y0 = t.Y*TileSize + TopMargin
	x1 = x0 + TileSize
	y1 = y0 + TileSize
	return // using named return parameters
}

func (t *Tile) removeWall(d int) {
	var mask uint
	mask = MASK & (^bits[d])
	t.walls &= mask
}

func (t *Tile) recursiveBacktracker() {
	// dirs := [4]int{0, 1, 2, 3}
	// rand.Shuffle(len(dirs), func(i, j int) { dirs[i], dirs[j] = dirs[j], dirs[i] })
	dirs := rand.Perm(4)
	for d := 0; d < 4; d++ {
		dir := dirs[d]
		tn := t.edges[dir]
		if tn != nil && tn.visited == false {
			t.removeWall(dir)
			tn.removeWall(oppdirs[dir])
			tn.visited = true
			tn.recursiveBacktracker()
		}
	}
}

// Layout the tile
func (t *Tile) Layout() {
	t.homeX = float64(LeftMargin + t.X*TileSize)
	t.homeY = float64(TopMargin + t.Y*TileSize)
}

// Update the tile state (transitions, user input)
func (t *Tile) Update() error {
	return nil
}

func (t *Tile) debugText(screen *ebiten.Image, str string) {
	bound, _ := font.BoundString(Acme.small, str)
	w := (bound.Max.X - bound.Min.X).Ceil()
	h := (bound.Max.Y - bound.Min.Y).Ceil()
	x, y := t.homeX-overSize, t.homeY-overSize
	tx := int(x) + (TileSize-w)/2
	ty := int(y) + (TileSize-h)/2 + h
	var c color.Color = BasicColors["Black"]
	// ebitenutil.DrawRect(screen, float64(tx), float64(ty), float64(w), float64(h), c)
	text.Draw(screen, str, Acme.small, tx, ty, c)
}

// Draw renders a Tile object
func (t *Tile) Draw(screen *ebiten.Image) {

	// scale, point translation, rotate, object translation

	op := &ebiten.DrawImageOptions{}

	// Reset RGB (not Alpha) forcibly
	// tilesheet already has black shapes
	{
		// reducing alpha leaves the endcaps doubled
		op.ColorM.Scale(0, 0, 0, 1)
		r := float64(t.G.colorWall.R) / 0xff
		g := float64(t.G.colorWall.G) / 0xff
		b := float64(t.G.colorWall.B) / 0xff
		op.ColorM.Translate(r, g, b, 0)
	}

	op.GeoM.Translate(t.homeX-overSize, t.homeY-overSize)

	screen.DrawImage(t.tileImage, op)

	if DebugMode {
		if t.Y != 0 {
			ebitenutil.DrawLine(screen,
				t.homeX,
				t.homeY,
				t.homeX+float64(TileSize),
				t.homeY,
				BasicColors["Black"])
		}
		if t.X != 0 {
			ebitenutil.DrawLine(screen,
				t.homeX,
				t.homeY,
				t.homeX,
				t.homeY+float64(TileSize),
				BasicColors["Black"])
		}
	}
	// t.debugText(gridImage, fmt.Sprintf("%04b", t.walls))
}
