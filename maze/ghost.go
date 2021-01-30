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
	scale := float64(TileSize / 2)
	dc.MoveTo(polyCoords[0]*2+scale, polyCoords[1]*2+scale)
	for i := 2; i < len(polyCoords); {
		x := polyCoords[i]
		i++
		y := polyCoords[i]
		i++
		dc.LineTo(x*2+scale, y*2+scale)
	}
	dc.ClosePath()
	dc.Fill()
	dc.SetRGB(1, 1, 1)
	dc.DrawCircle(scale-10, scale-2, 8)
	dc.DrawCircle(scale+10, scale-2, 8)
	dc.Fill()
	dc.SetRGB(0, 0, 0)
	switch dir {
	case 0:
		dc.DrawCircle(scale-10, scale-4, 4)
		dc.DrawCircle(scale+10, scale-4, 4)
	case 1:
		dc.DrawCircle(scale-10+2, scale-2, 4)
		dc.DrawCircle(scale+10+2, scale-2, 4)
	case 2:
		dc.DrawCircle(scale-10, scale+2, 4)
		dc.DrawCircle(scale+10, scale+2, 4)
	case 3:
		dc.DrawCircle(scale-10-2, scale-2, 4)
		dc.DrawCircle(scale+10-2, scale-2, 4)
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

// Update the state/position of the Ghost
func (g *Ghost) Update() error {

	if g.dest == nil {
		var dirfuncs [4]func(int) int
		if rand.Float64() < 0.5 {
			dirfuncs = [4]func(int) int{left, forward, right, opposite}
		} else {
			dirfuncs = [4]func(int) int{right, forward, left, opposite}
		}
		for i := 0; i < 4; i++ {
			dir := dirfuncs[i](g.facing)
			if !g.tile.IsWall(dir) {
				g.facing = dir
				g.dest = g.tile.Neighbour(dir)
				break
			}
		}
		if g.dest != nil {
			g.lerpstep = 0.01
			g.srcX, g.srcY = g.tile.Position()
			g.dstX, g.dstY = g.dest.Position()
		} else {
			println("ghost has no direction")
		}
	} else {
		if g.lerpstep >= 1 {
			g.tile = g.dest
			g.worldX, g.worldY = g.tile.Position()
			g.dest = nil
		} else {
			g.worldX = lerp(g.srcX, g.dstX, g.lerpstep)
			g.worldY = lerp(g.srcY, g.dstY, g.lerpstep)
			g.lerpstep += 0.01
		}
	}

	return nil
}

// Draw the Ghost
func (g *Ghost) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.worldX, g.worldY)
	op.GeoM.Translate(CameraX, CameraY)
	screen.DrawImage(g.ghostImages[g.facing], op)
}
