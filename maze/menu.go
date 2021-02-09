// Copyright ©️ 2021 oddstream.games

package maze

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Menu represents a game state.
type Menu struct {
	widgets []Widget
	input   *Input
}

// NewMenu creates and initializes a Menu/GameState object
func NewMenu() *Menu {
	i := NewInput()
	s := &Menu{input: i}

	s.widgets = []Widget{
		NewLabel("CAN YOU HERD KITTENS?", TheAcmeFonts.large),
		NewLabel("Move the yellow blob by clicking where you want it to go", TheAcmeFonts.normal),
		NewLabel("Build/demolish walls using the WASD keys", TheAcmeFonts.normal),
		NewLabel("Herd the kittens into the square in the middle", TheAcmeFonts.normal),
		NewTextButton("START", 200, 50, TheAcmeFonts.normal, func() {
			TheUserData.CompletedLevels = 0
			TheUserData.Save()
			GSM.Switch(NewCutscene())
		}, i),
	}

	if TheUserData.CompletedLevels > 0 {
		s.widgets = append(s.widgets,
			NewTextButton("CONTINUE", 200, 50, TheAcmeFonts.normal, func() { GSM.Switch(NewCutscene()) }, i),
		)
	}

	return s
}

// Layout implements ebiten.Game's Layout
func (s *Menu) Layout(outsideWidth, outsideHeight int) (int, int) {

	xCenter := outsideWidth / 2
	yPlaces := []int{} // golang gotcha: can't use len(s.widgets) to make an array
	slots := len(s.widgets) + 1
	for i := 0; i < slots; i++ {
		yPlaces = append(yPlaces, (outsideHeight/slots)*i)
	}

	for i, w := range s.widgets {
		w.SetPosition(xCenter, yPlaces[i+1])
	}

	return outsideWidth, outsideHeight
}

// Update updates the current game state.
func (s *Menu) Update() error {

	s.input.Update()

	return nil
}

// Draw draws the current GameState to the given screen
func (s *Menu) Draw(screen *ebiten.Image) {
	screen.Fill(colorBackground)

	// op := &ebiten.DrawImageOptions{}

	for _, d := range s.widgets {
		d.Draw(screen)
	}
}
