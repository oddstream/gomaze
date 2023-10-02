package maze

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Gameover represents a game state.
type Gameover struct {
	widgets []Widget
	input   *Input
}

// NewGameover creates and initializes a Gameover/GameState object
func NewGameover() *Gameover {

	i := NewInput()
	g := &Gameover{input: i}

	g.widgets = []Widget{
		NewLabel("GAME OVER", TheAcmeFonts.large),
		NewTextButton("START", 200, 50, TheAcmeFonts.normal, func() {
			TheUserData.CompletedLevels = 0
			TheUserData.Save()
			GSM.Switch(NewCutscene())
		}, i),
	}

	return g
}

// Layout implements ebiten.Game's Layout
func (g *Gameover) Layout(outsideWidth, outsideHeight int) (int, int) {

	xCenter := outsideWidth / 2
	yPlaces := []int{} // golang gotcha: can't use len(s.widgets) to make an array
	slots := len(g.widgets) + 1
	for i := 0; i < slots; i++ {
		yPlaces = append(yPlaces, (outsideHeight/slots)*i)
	}

	for i, w := range g.widgets {
		w.SetPosition(xCenter, yPlaces[i+1])
	}

	return outsideWidth, outsideHeight
}

// Update updates the current game state.
func (g *Gameover) Update() error {

	g.input.Update()

	return nil
}

// Draw draws the current GameState to the given screen
func (g *Gameover) Draw(screen *ebiten.Image) {
	screen.Fill(colorBackground)
	for _, d := range g.widgets {
		d.Draw(screen)
	}
}
