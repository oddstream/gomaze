// Copyright ©️ 2021 oddstream.games

package maze

import (
	"math/rand"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gomaze/util"
)

// Ghost defines the yellow blob/player avatar
type Ghost struct {
	tile                   *Tile   // tile we are sitting on
	dest                   *Tile   // tile we are lerping to
	facing                 int     // 0,1,2,3
	srcX, srcY, dstX, dstY float64 // positions for lerp
	lerpstep               float64

	Images             [4]*ebiten.Image
	directionlessImage *ebiten.Image

	worldX, worldY float64
}

var polyCoords = []float64{
	-12, 10, // bottom left
	-8, 8,
	-4, 12,
	0, 8,
	4, 12,
	8, 8,
	12, 10, // bottom right
	12, 0,
	// arch
	11, -5,
	10, -7,
	9, -8,
	8, -9,
	6, -10,

	0, -11,

	-6, -10,
	-8, -9,
	-9, -8,
	-10, -7,
	-11, -5,
	// end of arch
	-12, 0,
}

func createImage(dir int) *ebiten.Image {
	mid := float64(TileSize / 2)
	dc := gg.NewContext(TileSize, TileSize)

	// comment this out to just draw googly eyes, which is kinda fun
	// dc.SetColor(BasicColors["Silver"])
	// dc.MoveTo(polyCoords[0]*2+mid, polyCoords[1]*2+mid)
	// for i := 2; i < len(polyCoords); {
	// 	x := polyCoords[i]
	// 	i++
	// 	y := polyCoords[i]
	// 	i++
	// 	dc.LineTo(x*2+mid, y*2+mid)
	// }
	// dc.ClosePath()
	// dc.Fill()
	// end

	dc.SetRGB(1, 1, 1)
	dc.DrawCircle(mid-10, mid-2, 8)
	dc.DrawCircle(mid+10, mid-2, 8)
	dc.Fill()
	dc.SetRGB(0, 0, 0)
	switch dir {
	case -1: // kludge for directionless
		dc.DrawCircle(mid-10, mid-2, 4)
		dc.DrawCircle(mid+10, mid-2, 4)
	case 0: // NORTH
		dc.DrawCircle(mid-10, mid-4, 4)
		dc.DrawCircle(mid+10, mid-4, 4)
	case 1: // EAST
		dc.DrawCircle(mid-10+2, mid-2, 4)
		dc.DrawCircle(mid+10+2, mid-2, 4)
	case 2: //SOUTH
		dc.DrawCircle(mid-10, mid+2, 4)
		dc.DrawCircle(mid+10, mid+2, 4)
	case 3: // WEST
		dc.DrawCircle(mid-10-2, mid-2, 4)
		dc.DrawCircle(mid+10-2, mid-2, 4)
	}
	dc.Fill()
	dc.Stroke()
	return ebiten.NewImageFromImage(dc.Image())
}

// NewGhost creates a new Ghost object
func NewGhost(start *Tile) *Ghost {
	g := &Ghost{tile: start}

	for d := 0; d < 4; d++ {
		g.Images[d] = createImage(d)
	}
	g.directionlessImage = createImage(-1)

	g.facing = 0
	g.worldX, g.worldY = g.tile.Position()

	return g
}

func (gh *Ghost) isDirOkay(dir int) bool {
	// can't go through walls
	if gh.tile.IsWall(dir) {
		return false
	}
	// don't like going where puck is
	if TheGrid.puck.tile == gh.tile.Neighbour(dir) {
		return false
	}
	// don't leave the pen - this makes it too easy
	// if gh.tile.pen {
	// 	tn := gh.tile.Neighbour(dir)
	// 	if !tn.pen {
	// 		return false
	// 	}
	// }
	// don't like going on top of other ghosts
	for _, g := range TheGrid.ghosts {
		if g == gh {
			continue
		}
		if g.dest == gh.tile.Neighbour(dir) {
			return false
		}
	}
	return true
}

// Update the state/position of the Ghost
func (gh *Ghost) Update() error {

	if gh.dest == nil {
		var dirfuncs [4]func(int) int
		if rand.Float64() < 0.5 {
			dirfuncs = [4]func(int) int{util.Leftward, util.Forward, util.Rightward, util.Backward}
		} else {
			dirfuncs = [4]func(int) int{util.Rightward, util.Forward, util.Leftward, util.Backward}
		}
		for i := 0; i < 4; i++ {
			dir := dirfuncs[i](gh.facing)
			if gh.isDirOkay(dir) {
				gh.facing = dir
				gh.dest = gh.tile.Neighbour(dir)
				break
			}
		}
		if gh.dest != nil {
			gh.lerpstep = 0.01
			gh.srcX, gh.srcY = gh.tile.Position()
			gh.dstX, gh.dstY = gh.dest.Position()
		} else {
			// ghost has no direction
			// this can happen when trying to stop ghosts from sitting on top of each other in Ghost.IsGoodDir()
		}
	} else {
		if gh.lerpstep >= 1 {
			gh.tile = gh.dest
			gh.worldX, gh.worldY = gh.tile.Position()
			gh.dest = nil
		} else {
			gh.worldX = util.Lerp(gh.srcX, gh.dstX, gh.lerpstep)
			gh.worldY = util.Lerp(gh.srcY, gh.dstY, gh.lerpstep)
			gh.lerpstep += 0.01
		}
	}

	return nil
}

// Draw the Ghost
func (gh *Ghost) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(gh.worldX, gh.worldY)
	op.GeoM.Translate(CameraX, CameraY)
	if gh.dest == nil {
		screen.DrawImage(gh.directionlessImage, op)
	} else {
		screen.DrawImage(gh.Images[gh.facing], op)
	}
}
