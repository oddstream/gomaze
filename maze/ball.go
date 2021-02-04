// Copyright ©️ 2021 oddstream.games

package maze

import (
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gomaze/util"
)

// Ball defines the yellow blob/player avatar
type Ball struct {
	tile *Tile // Tile we are sitting on
	dest *Tile // Tile we have been thrown to

	ballImage *ebiten.Image

	srcX, srcY, dstX, dstY float64 // positions for lerp
	lerpstep               float64

	worldX, worldY float64
}

// NewBall creates a new Ball object
func NewBall(start *Tile) *Ball {
	b := &Ball{tile: start}

	dc := gg.NewContext(TileSize, TileSize)
	dc.SetColor(BasicColors["Yellow"])
	dc.DrawCircle(float64(TileSize/2), float64(TileSize/2), float64(TileSize/8))
	dc.Fill()
	dc.Stroke()
	b.ballImage = ebiten.NewImageFromImage(dc.Image())

	b.worldX, b.worldY = b.tile.Position()

	return b
}

// ThrowTo a target tile
func (b *Ball) ThrowTo(to *Tile) {
	b.dest = to
	b.srcX, b.srcY = b.tile.Position()
	b.dstX, b.dstY = b.dest.Position()
	b.lerpstep = 0.05
}

// Tile getter for ball's location
// func (b *Ball) Tile() *Tile {
// 	return b.tile
// }

// Update the state/position of the Ball
func (b *Ball) Update() error {
	// println("tile=", b.tile, "targ=", b.targ)
	if b.dest != nil {
		if b.lerpstep >= 1 {
			b.tile = b.dest
			b.dest = nil
			b.worldX, b.worldY = b.tile.Position()
		} else {
			b.worldX = util.Smoothstep(b.srcX, b.dstX, b.lerpstep)
			b.worldY = util.Smoothstep(b.srcY, b.dstY, b.lerpstep)
			b.lerpstep += 0.05
		}
	}
	return nil
}

// Draw the Ball
func (b *Ball) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.worldX, b.worldY)
	op.GeoM.Translate(CameraX, CameraY)
	screen.DrawImage(b.ballImage, op)
}
