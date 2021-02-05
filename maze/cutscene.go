// Copyright ©️ 2021 oddstream.games

package maze

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/fogleman/gg"
)

// Cutscene represents a game state.
type Cutscene struct {
	newWidth, newHeight, newGhosts int
	circleImage                    *ebiten.Image
	circlePos                      image.Point
	skew                           float64
}

// NewCutscene creates and initializes a Cutscene/GameState object
func NewCutscene(newWidth, newHeight, newGhosts int) *Cutscene {
	cs := &Cutscene{newWidth: newWidth, newHeight: newHeight, newGhosts: newGhosts}

	dc := gg.NewContext(400, 400)
	dc.SetRGB(1, 1, 0)
	dc.DrawCircle(200, 200, 120)
	dc.Fill()
	dc.Stroke()
	cs.circleImage = ebiten.NewImageFromImage(dc.Image())

	return cs
}

// Layout implements ebiten.Game's Layout
func (cs *Cutscene) Layout(outsideWidth, outsideHeight int) (int, int) {

	xCenter := outsideWidth / 2
	yCenter := outsideHeight / 2

	cx, cy := cs.circleImage.Size()
	cs.circlePos = image.Point{X: xCenter - (cx / 2), Y: yCenter - (cy / 2)}

	return outsideWidth, outsideHeight
}

// Update updates the current game state.
func (cs *Cutscene) Update() error {

	if cs.skew < 90 {
		cs.skew++
	} else {
		GSM.Switch(NewGrid(cs.newWidth, cs.newHeight, cs.newGhosts))
	}

	return nil
}

// Draw draws the current GameState to the given screen
func (cs *Cutscene) Draw(screen *ebiten.Image) {
	screen.Fill(BasicColors["Black"])

	skewRadians := cs.skew * math.Pi / 180

	{
		op := &ebiten.DrawImageOptions{}
		sx, sy := cs.circleImage.Size()
		sx, sy = sx/2, sy/2
		op.GeoM.Translate(float64(-sx), float64(-sy))
		op.GeoM.Scale(0.5, 0.5)
		op.GeoM.Skew(skewRadians, skewRadians)
		op.GeoM.Translate(float64(sx), float64(sy))

		op.GeoM.Translate(float64(cs.circlePos.X), float64(cs.circlePos.Y))
		screen.DrawImage(cs.circleImage, op)
	}
}
