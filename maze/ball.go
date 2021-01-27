// Copyright ©️ 2021 oddstream.games

package maze

import (
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

// Ball defines the yellow blob/player avatar
type Ball struct {
	tile *Tile // Tile we are sitting on
	targ *Tile // Tile we have been thrown to

	ballImage *ebiten.Image

	srcX, srcY, dstX, dstY float64 // positions for lerp
	lerpstep               float64

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
	b.srcX, b.srcY = b.tile.Position()
	b.dstX, b.dstY = b.targ.Position()
	b.lerpstep = 0.05
}

// Tile getter for ball's location
// func (b *Ball) Tile() *Tile {
// 	return b.tile
// }

// Update the state/position of the Ball
func (b *Ball) Update() error {
	// println("tile=", b.tile, "targ=", b.targ)
	if b.targ != nil && b.tile != b.targ {
		if b.lerpstep >= 1 {
			b.tile = b.targ
			b.targ = nil
			b.x, b.y = b.tile.Position()
		} else {
			b.x = smoothstep(b.srcX, b.dstX, b.lerpstep)
			b.y = smoothstep(b.srcY, b.dstY, b.lerpstep)
			b.lerpstep += 0.05
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
