// Copyright ©️ 2021 oddstream.games

package maze

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/fogleman/gg"
)

// Cutscene represents a game state.
type Cutscene struct {
	circleImage     *ebiten.Image
	circlePos       image.Point
	startX, finishX int
}

// NewCutscene creates and initializes a Cutscene/GameState object
func NewCutscene() *Cutscene {
	cs := &Cutscene{}

	dc := gg.NewContext(200, 200)
	dc.SetRGB(1, 1, 0)
	dc.DrawCircle(100, 100, 100)
	dc.Fill()
	dc.Stroke()
	cs.circleImage = ebiten.NewImageFromImage(dc.Image())

	return cs
}

// Layout implements ebiten.Game's Layout
func (cs *Cutscene) Layout(outsideWidth, outsideHeight int) (int, int) {

	if cs.circlePos.X == 0 && cs.circlePos.Y == 0 {
		// cx, cy := cs.circleImage.Size()
		cx := cs.circleImage.Bounds().Dx()
		cy := cs.circleImage.Bounds().Dy()
		cs.startX = -cx
		cs.finishX = outsideWidth
		cs.circlePos = image.Point{X: cs.startX, Y: outsideHeight - cy}
	}
	return outsideWidth, outsideHeight
}

// Update updates the current game state.
func (cs *Cutscene) Update() error {

	cs.circlePos.X += 20
	if cs.circlePos.X > cs.finishX {
		lvl := TheUserData.CompletedLevels
		TheGrid = NewGrid(LevelData[lvl][0], LevelData[lvl][1], LevelData[lvl][2])
		GSM.Switch(TheGrid)
	}

	return nil
}

// Draw draws the current GameState to the given screen
func (cs *Cutscene) Draw(screen *ebiten.Image) {
	screen.Fill(colorBackground)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(cs.circlePos.X), float64(cs.circlePos.Y))
	screen.DrawImage(cs.circleImage, op)
}
