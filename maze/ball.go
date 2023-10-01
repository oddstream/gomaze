// Copyright ©️ 2021 oddstream.games

package maze

import (
	"fmt"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gomaze/util"
)

// Ball defines the target destination of the yellow blob/player avatar
type Ball struct {
	tile                   *Tile // Tile ball is sitting on
	dest                   *Tile // Tile ball has been thrown to
	ballImage              *ebiten.Image
	srcX, srcY, dstX, dstY float64 // positions for lerp
	lerpstep               float64 // lerp 0.0 .. 1.0
	worldX, worldY         float64
}

// NewBall creates a new Ball object
func NewBall(start *Tile, col color.RGBA) *Ball {
	b := &Ball{tile: start}

	dc := gg.NewContext(TileSize, TileSize)
	dc.SetColor(col)
	dc.DrawCircle(float64(TileSize/2), float64(TileSize/2), float64(TileSize/8))
	dc.Fill()
	dc.Stroke()
	b.ballImage = ebiten.NewImageFromImage(dc.Image())

	b.worldX, b.worldY = b.tile.position()

	return b
}

// StartThrow to a target tile
func (b *Ball) StartThrow(to *Tile) {
	b.dest = to
	b.srcX, b.srcY = b.tile.position()
	b.dstX, b.dstY = b.dest.position()
	b.lerpstep = 0.025
}

// String representation of this ball
func (b *Ball) String() string {
	return fmt.Sprintf("(%v,%v)", b.tile.X, b.tile.Y)
}

// Tile getter for ball's location
// func (b *Ball) Tile() *Tile {
// 	return b.tile
// }

// Update the state/position of the Ball
func (b *Ball) Update() error {
	// println("tile=", b.tile, "targ=", b.targ)
	if b.dest == nil {
		return nil
	}
	if b.lerpstep >= 1.0 {
		b.tile = b.dest
		b.dest = nil
		b.worldX, b.worldY = b.tile.position()
	} else {
		b.worldX = util.Smoothstep(b.srcX, b.dstX, b.lerpstep)
		b.worldY = util.Smoothstep(b.srcY, b.dstY, b.lerpstep)
		b.lerpstep += 0.025
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
