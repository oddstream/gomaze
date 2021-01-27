// Copyright ©️ 2021 oddstream.games

package maze

import (
	"math"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

// Ball defines the yellow blob/player avatar
type Ball struct {
	tile *Tile // Tile we are sitting on
	targ *Tile // Tile we have been thrown to

	ballImage *ebiten.Image

	x, y float64
}

// NewBall creates a new Ball object
func NewBall(start *Tile) *Ball {
	b := &Ball{tile: start, targ: start}

	dc := gg.NewContext(TileSize, TileSize)
	dc.SetRGB(1, 1, 0)
	dc.DrawCircle(float64(TileSize/2), float64(TileSize/2), float64(TileSize/8))
	dc.Fill()
	dc.Stroke()
	b.ballImage = ebiten.NewImageFromImage(dc.Image())

	b.x, b.y = b.tile.Position()

	return b
}

// ThrowTo a target tile
func (b *Ball) ThrowTo(t *Tile) {
	b.targ = t
}

// Tile getter for ball's location
// func (b *Ball) Tile() *Tile {
// 	return b.tile
// }

// Update the state/position of the Ball
func (b *Ball) Update() error {
	// println("tile=", b.tile, "targ=", b.targ)
	if b.tile != b.targ {
		dx, dy := b.targ.Position()
		// println("ball moving from", b.x, b.y, "to", dx, dy)
		b.x = lerp(b.x, dx, 0.1)
		b.y = lerp(b.y, dy, 0.1)
		if math.Abs(dx-b.x) < float64(TileSize/4) && math.Abs(dy-b.y) < float64(TileSize/4) {
			b.tile = b.targ
			b.x, b.y = b.tile.Position()
		}
	}
	return nil
}

// Draw the Ball
func (b *Ball) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.x, b.y)
	screen.DrawImage(b.ballImage, op)
}
