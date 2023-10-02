package maze

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"runtime"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"oddstream.games/gomaze/util"
)

const (
	// TileSize is now a constant
	TileSize = 80
	// MaxGhosts limit that can fit in 3x3 pen
	MaxGhosts = 8
)

// TilesAcross and TilesDown are package-level variables so they can be seen by Tile
var (
	TilesAcross int
	TilesDown   int
	CameraX     float64
	CameraY     float64
)

// Grid is an object representing the grid of tiles
type Grid struct {
	ticks              int     // count update ticks (use ticks % 10 for refreshing minimap)
	tiles              []*Tile // a slice (not array!) of pointers to Tile objects
	colorBackground    color.RGBA
	colorWall          color.RGBA
	input              *Input
	puck               *Puck
	ghosts             []*Ghost
	meanies            []*Meanie
	penImage           *ebiten.Image
	penX, penY         float64
	minimapImage       *ebiten.Image
	minimapX, minimapY float64
}

// NewGrid create a Grid object
func NewGrid(w, h, ghostCount int) *Grid {

	// var screenWidth, screenHeight int

	// if runtime.GOARCH == "wasm" {
	// 	screenWidth, screenHeight = WindowWidth, WindowHeight
	// } else {
	// 	screenWidth, screenHeight = ebiten.WindowSize()
	// }

	TilesAcross, TilesDown = w, h

	g := &Grid{tiles: make([]*Tile, TilesAcross*TilesDown)}
	for i := range g.tiles {
		g.tiles[i] = NewTile(i%TilesAcross, i/TilesAcross)
	}

	// link the tiles together to avoid all that tedious 2d array stuff
	for _, t := range g.tiles {
		x := t.X
		y := t.Y
		t.edges[0] = g.findTile(x, y-1) // North
		t.edges[1] = g.findTile(x+1, y) // East
		t.edges[2] = g.findTile(x, y+1) // South
		t.edges[3] = g.findTile(x-1, y) // West
	}

	{
		// the pen is a 3x3 square of tiles in the middle of the grid
		midX := TilesAcross / 2
		midY := TilesDown / 2
		for x := midX - 1; x <= midX+1; x++ {
			for y := midY - 1; y <= midY+1; y++ {
				t := g.findTile(x, y)
				t.removeAllWalls()
				t.pen = true
			}
		}
		// create a shaded background image for the pen
		dc := gg.NewContext(TileSize*3, TileSize*3)
		dc.SetRGBA(float64(g.colorBackground.R/0xff), float64(g.colorBackground.G/0xff), float64(g.colorBackground.B/0xff), 0.5)
		dc.DrawRoundedRectangle(0, 0, float64(TileSize*3), float64(TileSize*3), float64(TileSize/12))
		dc.Fill()
		dc.Stroke()

		// put the level number dimly in the background image
		dc.SetRGBA(float64(g.colorBackground.R/0xff), float64(g.colorBackground.G/0xff), float64(g.colorBackground.B/0xff), 0.1)
		dc.SetFontFace(TheAcmeFonts.huge)
		dc.DrawStringAnchored(fmt.Sprint(TheUserData.CompletedLevels+1), float64(TileSize)*1.5, float64(TileSize), 0.5, 0.5)

		g.penImage = ebiten.NewImageFromImage(dc.Image())

		// remember where to draw the pen background image
		t := g.findTile(midX-1, midY-1)
		x, y, _, _ := t.rect()
		g.penX = float64(x)
		g.penY = float64(y)
	}

	g.carve()

	g.puck = NewPuck(g.findTile(TilesAcross/2, TilesDown/2), BasicColors["Yellow"])

	if ghostCount > MaxGhosts {
		ghostCount = MaxGhosts
	}
	for i := 0; i < ghostCount; i++ {
		g.ghosts = append(g.ghosts, NewGhost(g.randomTile()))
	}

	palette := Palettes[rand.Int()%len(Palettes)]
	g.colorBackground = CalcBackgroundColor(palette)
	g.colorWall = ExtendedColors[palette[rand.Int()%len(palette)]]

	g.input = NewInput()
	g.input.Add(g)

	return g
}

// NotifyCallback is called by the Subject (Input) when something interesting happens
func (g *Grid) NotifyCallback(event interface{}) {
	switch v := event.(type) { // Type switch https://tour.golang.org/methods/16
	case image.Point:
		pt := v
		pt.X = pt.X - int(CameraX)
		pt.Y = pt.Y - int(CameraY)
		// pt = pt.Sub(image.Point{X: int(CameraX), Y: int(CameraY)})
		t := g.findTileAt(pt)
		if t != nil {
			if t == g.puck.tile {
				// println("Puck is already on tile", t.X, t.Y, t.WhichQuadrant(pt))
				dir := t.whichQuadrant(pt)
				if t.neighbour(dir) != nil {
					if t.isWall(dir) {
						t.removeWall(dir)
					} else {
						t.addWall(dir)
					}
				}
				// g.puck.tile.toggleWall(t.whichQuadrant(pt))
				g.visitTiles()
			} else {
				// println("input on tile", t.X, t.Y, t.wallCount())
				g.AllTiles(func(t *Tile) { t.parent = nil; t.marked = false })
				g.puck.ThrowBallTo(t)
			}
		}
	case ebiten.Key:
		k := v
		switch k {
		case ebiten.KeyBackspace:
			GSM.Switch(NewMenu())
		case ebiten.KeyW:
			g.puck.tile.toggleWall(0)
			g.visitTiles()
		case ebiten.KeyD:
			g.puck.tile.toggleWall(1)
			g.visitTiles()
		case ebiten.KeyS:
			g.puck.tile.toggleWall(2)
			g.visitTiles()
		case ebiten.KeyA:
			g.puck.tile.toggleWall(3)
			g.visitTiles()
		case ebiten.KeyM:
			g.meanies = append(g.meanies, NewMeanie(g.randomTile()))
		}
	}
}

// Size returns the size of the grid in pixels
func (g *Grid) Size() (int, int) {
	return TilesAcross * TileSize, TilesDown * TileSize
}

func (g *Grid) findTile(x, y int) *Tile {

	// we re-order the tiles when dragging, to put the dragged tile at the top of the z-order
	// so we can't use i := x + (y * TilesAcross) to find index of tile in slice
	// except if this func is just used after tiles have been created

	// for _, t := range g.tiles {
	// 	if t.X == x && t.Y == y {
	// 		return t
	// 	}
	// }
	// return nil

	if x < 0 || x >= TilesAcross {
		return nil
	}
	if y < 0 || y >= TilesDown {
		return nil
	}
	i := x + (y * TilesAcross)
	if i < 0 || i > len(g.tiles) {
		log.Fatal("findTile index out of bounds")
	}
	return g.tiles[i]
}

// findTileAt finds the tile under the mouse click or touch
func (g *Grid) findTileAt(pt image.Point) *Tile {
	for _, t := range g.tiles {
		if util.InRect(pt, t.rect) {
			return t
		}
	}
	return nil
}

func (g *Grid) randomTile() *Tile {
	i := rand.Intn(len(g.tiles))
	return g.tiles[i]
}

// BFS does a breadth first search and sets Tile.parent
// to form a path from start to dest tiles.
func (g *Grid) BFS(start, dst *Tile) {
	g.AllTiles(func(t *Tile) { t.parent = nil })
	q := []*Tile{start}
	start.parent = start
	for len(q) > 0 {
		t := q[0]
		q = q[1:] // take first tile off front of queue
		if t == dst {
			return
		}
		for _, d := range ALL_DIRECTIONS {
			tn := t.neighbour(d)
			if tn != nil && !t.isWall(d) && tn.parent == nil {
				tn.parent = t
				q = append(q, tn)
			}
		}
	}
	// if puck has walled a ghost in, then a path will not be found
	// log.Println("Grid.BFS: path not found")
}

func (g *Grid) findTileTowards(src, dst *Tile) *Tile {
	g.BFS(src, dst)
	for dst != nil && dst.parent != src {
		dst = dst.parent
	}
	return dst
}

// func (g *Grid) removeRandomWall() {
// 	t := g.randomTile()
// 	for _, d := range []int{0, 1, 2, 3} {
// 		if t.isWall(d) && t.neighbour(d) != nil {
// 			t.removeWall(d)
// 			break
// 		}
// 	}
// }

// func (g *Grid) createRoom(x1, y1 int) {
// 	x2 := x1 + 4
// 	y2 := y1 + 4
// 	if x2 >= TilesAcross || y2 >= TilesDown {
// 		return
// 	}
// 	for x := x1; x < x2; x++ {
// 		for y := y1; y < y2; y++ {
// 			t := g.findTile(x, y)
// 			t.removeAllWalls()
// 		}
// 	}
// 	for x := x1; x < x2; x++ {
// 		t := g.findTile(x, y1)
// 		t.addWall(0)
// 		t = g.findTile(x, y2-1)
// 		t.addWall(2)
// 	}
// 	for y := y1; y < y2; y++ {
// 		t := g.findTile(x1, y)
// 		t.addWall(3)
// 		t = g.findTile(x2-1, y)
// 		t.addWall(1)
// 	}
// }

// func (g *Grid) createRooms() {
// 	for x := 0; x < TilesAcross; x += 4 {
// 		for y := 0; y < TilesDown; y += 4 {
// 			if rand.Float64() < 0.5 {
// 				g.createRoom(x, y)
// 			}
// 		}
// 	}
// }

func (g *Grid) carve() {
	t := g.randomTile()
	// t.recursiveBacktracker()
	t.prim()
}

// func (g *Grid) fillCulDeSacs() {
// 	for _, t := range g.tiles {
// 		t.fillCulDeSac()
// 	}
// }

// AllTiles applies a func to all tiles
func (g *Grid) AllTiles(fn func(*Tile)) {
	for _, t := range g.tiles {
		fn(t)
	}
}

// visitTiles rewalks the grid (from the middle) and resets the .visited
// member, which is used to differentiate reachable and unreachable tiles;
// .visited == true means tile is reachable
func (g *Grid) visitTiles() {
	g.AllTiles(func(t *Tile) { t.visited = false; t.parent = nil })
	tFirst := g.findTile(TilesAcross/2, TilesDown/2)
	q := []*Tile{tFirst}
	tFirst.parent = tFirst
	for len(q) > 0 {
		t := q[0]
		t.visited = true
		q = q[1:] // take first tile off front of queue
		for _, d := range ALL_DIRECTIONS {
			if t.isWall(d) {
				continue
			}
			tn := t.neighbour(d)
			if tn == nil {
				log.Fatal("open unwalled edge found in visitTiles BFS")
			}
			if tn.parent == nil {
				tn.parent = t
				q = append(q, tn)
			}
		}
	}
}

// getMinimap shows position of ghosts
func (g *Grid) getMinimap(screen *ebiten.Image) *ebiten.Image {

	worldWidth, worldHeight := float64(TilesAcross*TileSize), float64(TilesDown*TileSize)
	mapWidth, mapHeight := worldWidth/10, worldHeight/10

	if g.ticks%10 != 0 && g.minimapImage != nil {
		return g.minimapImage
	}

	halfTileSize := float64(TileSize / 2)

	dc := gg.NewContext(int(mapWidth), int(mapHeight))
	dc.DrawRectangle(0, 0, float64(mapWidth-1), float64(mapHeight-1))
	dc.SetRGB(1, 1, 1)
	dc.Stroke()

	for _, gh := range g.ghosts {
		x := util.MapValue(gh.worldX+halfTileSize, 0, worldWidth, 0, mapWidth)
		y := util.MapValue(gh.worldY+halfTileSize, 0, worldHeight, 0, mapHeight)
		dc.DrawCircle(x, y, 1)
	}
	dc.SetRGB(1, 1, 1)
	dc.Fill()

	for _, m := range g.meanies {
		x := util.MapValue(m.worldX+halfTileSize, 0, worldWidth, 0, mapWidth)
		y := util.MapValue(m.worldY+halfTileSize, 0, worldHeight, 0, mapHeight)
		dc.DrawCircle(x, y, 1)
	}
	dc.SetRGB(0, 0, 0)
	dc.Fill()

	{
		x := util.MapValue(g.puck.worldX+halfTileSize, 0, worldWidth, 0, mapWidth)
		y := util.MapValue(g.puck.worldY+halfTileSize, 0, worldHeight, 0, mapHeight)
		dc.DrawCircle(x, y, 2)
		dc.SetRGB(1, 1, 0)
		dc.Fill()
	}

	dc.Stroke()

	g.minimapImage = ebiten.NewImageFromImage(dc.Image())

	return g.minimapImage
}

// Layout implements ebiten.Game's Layout.
func (g *Grid) Layout(outsideWidth, outsideHeight int) (int, int) {

	worldWidth := float64(TilesAcross * TileSize)
	mapWidth := worldWidth / 10
	g.minimapX, g.minimapY = float64(outsideWidth)-mapWidth, 0

	for _, t := range g.tiles {
		t.Layout()
	}
	// g.puck.Layout(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}

// Update the board state (transitions, user input)
func (g *Grid) Update() error {

	g.ticks++
	if g.ticks > int(ebiten.ActualTPS()) {
		g.ticks = 0
	}

	g.input.Update()

	for _, t := range g.tiles {
		t.Update() // doesn't do anything at the moment
	}

	count := 0
	for _, gh := range g.ghosts {
		gh.Update()
		if gh.tile.pen {
			count++
		}
	}
	if count == len(g.ghosts) {
		TheUserData.CompletedLevels++
		if TheUserData.CompletedLevels >= len(LevelData) {
			TheUserData.CompletedLevels = 0
			TheUserData.Save()
			GSM.Switch(NewGameover())
		} else {
			TheUserData.Save()
			GSM.Switch(NewCutscene())
		}
	}
	for _, m := range g.meanies {
		m.Update()
	}

	g.puck.Update()

	return nil
}

// Draw renders the grid into the gridImage
func (g *Grid) Draw(screen *ebiten.Image) {

	screen.Fill(g.colorBackground)

	{
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(g.penX+CameraX, g.penY+CameraY)
		// op.GeoM.Translate(CameraX, CameraY)
		screen.DrawImage(g.penImage, op)
	}

	for _, t := range g.tiles {
		t.Draw(screen)
	}

	for _, t := range g.tiles {
		t.DrawMarked(screen)
	}

	for _, gh := range g.ghosts {
		gh.Draw(screen)
	}

	for _, m := range g.meanies {
		m.Draw(screen)
	}

	g.puck.Draw(screen)

	{
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(g.minimapX, g.minimapY)
		screen.DrawImage(g.getMinimap(screen), op)
	}

	if DebugMode {
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		ebitenutil.DebugPrint(screen, fmt.Sprintf("NumGC %v", ms.NumGC))
		// ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS %v, FPS %v, grid %d,%d camera %v,%v puck %v,%v",
		// 	math.Ceil(ebiten.CurrentTPS()),
		// 	math.Ceil(ebiten.CurrentFPS()),
		// 	TilesAcross, TilesDown,
		// 	CameraX, CameraY,
		// 	g.puck.tile.X, g.puck.tile.Y))
	}
}
