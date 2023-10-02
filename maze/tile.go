package maze

import (
	"fmt"
	"image"
	"log"
	"math/rand"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gomaze/util"
)

// NORTH_WALL EAST_WALL SOUTH_WALL WEST_WALL bit patterns for presence of walls
const (
	NORTH_WALL = 0b0001 // 1 << iota
	EAST_WALL  = 0b0010 // 1 << 1
	SOUTH_WALL = 0b0100 // 1 << 2
	WEST_WALL  = 0b1000 // 1 << 3
	ALL_WALLS  = 0b1111

	NORTH = 0
	EAST  = 1
	SOUTH = 2
	WEST  = 3
)

var (
	reachableImages   map[uint]*ebiten.Image
	unreachableImages map[uint]*ebiten.Image
	dotImage          *ebiten.Image
	overSize          float64
	wallbits          = [4]uint{NORTH_WALL, EAST_WALL, SOUTH_WALL, WEST_WALL} // map a direction (0..3) to it's bits
	wallopps          = [4]uint{SOUTH_WALL, WEST_WALL, NORTH_WALL, EAST_WALL} // map a direction (0..3) to it's opposite bits
	// oppdirs  = [4]int{2, 3, 0, 1}
	ALL_DIRECTIONS = [4]Direction{NORTH, EAST, SOUTH, WEST}
)

type Direction int

func init() {
	if TileSize == 0 {
		log.Fatal("Tile dimensions not set")
	}

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
	actualTileSize := reachableImages[0].Bounds().Dx()
	overSize = float64((actualTileSize - TileSize) / 2)

	{
		mid := float64(actualTileSize / 2)
		dc := gg.NewContext(actualTileSize, actualTileSize)
		dc.SetRGBA(0, 0, 0, 0.5)
		dc.DrawCircle(mid, mid, 2)
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
	pen            bool // true if this tile is part of the central pen

	// members that may change
	walls uint // bit mask of the walls this tile currently has

	// volatile members
	marked  bool  // true if this tile is part of marked path for puck to follow
	visited bool  // temporarily used when doing BFS of grid
	parent  *Tile // temporarily used to find puck's marked path
}

// NewTile creates a new Tile object and returns a pointer to it
// all new tiles start with all four walls before they are carved later
func NewTile(x, y int) *Tile {
	t := &Tile{X: x, Y: y, walls: ALL_WALLS}
	// worldX, worldY will be (re)set by Layout()
	return t
}

// reset prepares a Tile for a new level by resetting just gameplay data, not structural data
// func (t *Tile) reset() {
// 	t.walls = ALL_WALLS
// 	t.visited = false
// 	t.parent = nil
// }

// rect gives the x,y screen coords of the tile's top left and bottom right corners
func (t *Tile) rect() (x0 int, y0 int, x1 int, y1 int) {
	x0 = t.X * TileSize
	y0 = t.Y * TileSize
	x1 = x0 + TileSize
	y1 = y0 + TileSize
	return // using named return parameters
}

// neighbour returns the neighbouring tile in that direction
func (t *Tile) neighbour(d Direction) *Tile {
	return t.edges[d]
}

// isWall returns true if there is a wall in that direction
func (t *Tile) isWall(d Direction) bool {
	bit := wallbits[d]
	return t.walls&bit == bit
}

func (t *Tile) addWall(d Direction) {
	t.walls |= wallbits[d]
	if tn := t.neighbour(d); tn != nil {
		tn.walls |= wallopps[d]
	}
}

func (t *Tile) removeWall(d Direction) {
	if tn := t.neighbour(d); tn != nil {
		var mask uint
		// unset the wall bit on this tile
		mask = ALL_WALLS & ^wallbits[d]
		t.walls &= mask
		// unset the opposite wall bit on the neighbour tile
		mask = ALL_WALLS & ^wallopps[d]
		tn.walls &= mask
	}
}

func (t *Tile) toggleWall(d Direction) {
	if tn := t.neighbour(d); tn != nil {
		if t.isWall(d) {
			t.removeWall(d)
		} else {
			t.addWall(d)
		}
	}
}

func (t *Tile) removeAllWalls() {
	for _, d := range ALL_DIRECTIONS {
		t.removeWall(d)
	}
}

// func (t *Tile) wallCount() int {
// 	return bits.OnesCount(t.walls)
// }

// func (t *Tile) fillCulDeSac() {
// 	if t.wallCount() == 3 {
// 		for d := 0; d < 4; d++ {
// 			if t.walls&wallbits[d] == 0 {
// 				t.walls = MASK
// 				if tn := t.Neighbour(d); tn != nil {
// 					tn.walls |= wallopps[d]
// 				}
// 				break
// 			}
// 		}
// 	}
// }

func (t *Tile) recursiveBacktracker() {
	dirs := [4]Direction{0, 1, 2, 3}
	rand.Shuffle(len(dirs), func(i, j int) { dirs[i], dirs[j] = dirs[j], dirs[i] })
	// dirs := rand.Perm(4)
	for d := 0; d < 4; d++ {
		dir := dirs[d]
		tn := t.neighbour(dir)
		if tn != nil && !tn.visited {
			t.removeWall(dir)
			tn.visited = true
			tn.recursiveBacktracker()
		}
	}
}

// mark this tile as 'in'
// and then adds marked/in neighbours to the frontier tiles
func (t *Tile) primMark(pfrontier *[]*Tile) {
	t.visited = true
	for _, dir := range ALL_DIRECTIONS {
		n := t.edges[dir]
		if n != nil && !n.visited {
			*pfrontier = append(*pfrontier, n)
		}
	}
}

// return all the 'in' neighbours
func (t *Tile) primNeighbours() []*Tile {
	var lst []*Tile = []*Tile{}
	for _, dir := range ALL_DIRECTIONS {
		n := t.edges[dir]
		if n != nil && n.visited {
			lst = append(lst, n)
		}
	}
	return lst
}

func primWhichDirIs(src, dst *Tile) Direction {
	for _, dir := range ALL_DIRECTIONS {
		if src.edges[dir] == dst {
			return dir
		}
	}
	return -1
}

func (t *Tile) prim() {
	if t.visited {
		log.Fatal("Tile.visited is true")
	}
	var frontier []*Tile = []*Tile{}
	t.primMark(&frontier)
	for len(frontier) > 0 {
		// remove a random frontier tile

		// we do not care about ordering,
		// so replace the element to delete with the one at the end of the slice
		// and then return the first len-1 elements
		i := rand.Intn(len(frontier))
		t1 := frontier[i]
		frontier[i] = frontier[len(frontier)-1]
		frontier = frontier[:len(frontier)-1]
		lst := t1.primNeighbours()
		t2 := lst[rand.Intn(len(lst))]
		if (t1.visited && !t2.visited) || (!t1.visited && t2.visited) {
			// remove wall between t1 and t2
			dir := primWhichDirIs(t1, t2)
			if dir == -1 {
				log.Fatal("dir is -1")
			}
			t1.removeWall(dir)
		}
		// mark the frontier tile as being 'in' the maze
		// (and add any of it's outside neighbours to the frontier)
		t1.primMark(&frontier)
	}
}

// position of this tile (top left origin) in screen coords
func (t *Tile) position() (float64, float64) {
	// Tile.Layout() may not have been called yet
	return float64(t.X * TileSize), float64(t.Y * TileSize)
	// return t.worldX, t.worldY
}

// whichQuadrant - given a point (in screen/world coords) that is assumed to be on this tile,
// return which direction/quadrant (0,1,2,3) the point is within. Imagine a diagonal cross
// on the tile. TODO this looks fugly.
func (t *Tile) whichQuadrant(pt image.Point) Direction {
	topLeft := image.Point{X: int(t.worldX), Y: int(t.worldY)}
	topRight := image.Point{X: int(t.worldX) + TileSize, Y: int(t.worldY)}
	center := image.Point{X: int(t.worldX) + (TileSize / 2), Y: int(t.worldY) + (TileSize / 2)}
	bottomLeft := image.Point{X: int(t.worldX), Y: int(t.worldY) + TileSize}
	bottomRight := image.Point{X: int(t.worldX) + TileSize, Y: int(t.worldY) + TileSize}
	if util.PointInTriangle(pt, topLeft, topRight, center) {
		return NORTH
	}
	if util.PointInTriangle(pt, topRight, bottomRight, center) {
		return EAST
	}
	if util.PointInTriangle(pt, bottomLeft, bottomRight, center) {
		return SOUTH
	}
	if util.PointInTriangle(pt, topLeft, bottomLeft, center) {
		return WEST
	}

	return -1
}

// String representation of this tile
func (t *Tile) String() string {
	return fmt.Sprintf("[%v,%v]", t.X, t.Y)
}

// Layout the tile
func (t *Tile) Layout() {
	t.worldX = float64(t.X * TileSize)
	t.worldY = float64(t.Y * TileSize)
}

// Update the tile state (transitions, user input)
func (t *Tile) Update() error {
	return nil
}

// func (t *Tile) debugText(screen *ebiten.Image, str string) {
// 	bound, _ := font.BoundString(TheAcmeFonts.large, str)
// 	w := (bound.Max.X - bound.Min.X).Ceil()
// 	h := (bound.Max.Y - bound.Min.Y).Ceil()
// 	x, y := t.worldX-overSize, t.worldY-overSize
// 	tx := int(x) + (TileSize-w)/2
// 	ty := int(y) + (TileSize-h)/2 + h
// 	c := color.RGBA{R: 0xff - colorBackground.R, G: 0xff - colorBackground.G, B: 0xff - colorBackground.B, A: 0xff}
// 	// var c color.Color = BasicColors["Black"]
// 	// ebitenutil.DrawRect(screen, float64(tx), float64(ty), float64(w), float64(h), c)
// 	text.Draw(screen, str, TheAcmeFonts.large, tx, ty, c)
// }

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
	// tilesheet already has black wall shapes
	// turn black into TheGrid.colorWall
	{
		var r, g, b float64
		r = float64(TheGrid.colorWall.R) / 0xff
		g = float64(TheGrid.colorWall.G) / 0xff
		b = float64(TheGrid.colorWall.B) / 0xff
		op.ColorM.Translate(r, g, b, 0)
		// op.ColorScale.Scale(r, g, b, 1)
	}

	if t.visited {
		screen.DrawImage(reachableImages[t.walls], op)
	} else {
		screen.DrawImage(unreachableImages[t.walls], op)
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

// DrawMarked renders a mark on a Tile object
func (t *Tile) DrawMarked(screen *ebiten.Image) {
	// https://ebiten.org/documents/performancetips.html
	// batch drawing of similar objects: don't intermingle drawing of tileImage and dotImage objects
	if t.marked {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(t.worldX-overSize, t.worldY-overSize)
		op.GeoM.Translate(CameraX, CameraY)
		screen.DrawImage(dotImage, op)
	}
}
