// Copyright ©️ 2021 oddstream.games

package maze

import (
	"math/rand"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

// Ghost defines the yellow blob/player avatar
type Ghost struct {
	tile                   *Tile   // tile we are sitting on
	dest                   *Tile   // tile we are lerping to
	facing                 int     // 0,1,2,3
	srcX, srcY, dstX, dstY float64 // positions for lerp
	lerpstep               float64

	ghostImages [4]*ebiten.Image

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

func drawGhost(dir int) *ebiten.Image {
	dc := gg.NewContext(TileSize, TileSize)
	dc.SetRGB(0.8, 0.8, 0.9)
	mid := float64(TileSize / 2)
	dc.MoveTo(polyCoords[0]*2+mid, polyCoords[1]*2+mid)
	for i := 2; i < len(polyCoords); {
		x := polyCoords[i]
		i++
		y := polyCoords[i]
		i++
		dc.LineTo(x*2+mid, y*2+mid)
	}
	dc.ClosePath()
	dc.Fill()
	dc.SetRGB(1, 1, 1)
	dc.DrawCircle(mid-10, mid-2, 8)
	dc.DrawCircle(mid+10, mid-2, 8)
	dc.Fill()
	dc.SetRGB(0, 0, 0)
	switch dir {
	case 0:
		dc.DrawCircle(mid-10, mid-4, 4)
		dc.DrawCircle(mid+10, mid-4, 4)
	case 1:
		dc.DrawCircle(mid-10+2, mid-2, 4)
		dc.DrawCircle(mid+10+2, mid-2, 4)
	case 2:
		dc.DrawCircle(mid-10, mid+2, 4)
		dc.DrawCircle(mid+10, mid+2, 4)
	case 3:
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
		g.ghostImages[d] = drawGhost(d)
	}

	g.facing = 0
	g.worldX, g.worldY = g.tile.Position()

	return g
}

func (gh *Ghost) isDirOkay(dir int) bool {
	if gh.tile.IsWall(dir) {
		return false
	}
	if TheGrid.puck.tile == gh.tile.Neighbour(dir) {
		return false
	}
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
			dirfuncs = [4]func(int) int{left, forward, right, backward}
		} else {
			dirfuncs = [4]func(int) int{right, forward, left, backward}
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
			println("ghost has no direction")
		}
	} else {
		if gh.lerpstep >= 1 {
			gh.tile = gh.dest
			gh.worldX, gh.worldY = gh.tile.Position()
			gh.dest = nil
		} else {
			gh.worldX = lerp(gh.srcX, gh.dstX, gh.lerpstep)
			gh.worldY = lerp(gh.srcY, gh.dstY, gh.lerpstep)
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
	screen.DrawImage(gh.ghostImages[gh.facing], op)
}
