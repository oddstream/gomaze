// Copyright ©️ 2021 oddstream.games

package maze

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Widget type implements UpDate, Draw and Pushed
type Widget interface {
	Update() error
	Draw(*ebiten.Image)
	SetPosition(int, int)
	Rect() (int, int, int, int)
	Pushed(*Input) bool
	Action()
}

// Pushable type implements Rect
// type Pushable interface {
// 	Rect() (int, int, int, int)
// }

// Menu represents a game state.
type Menu struct {
	widgets []Widget
	input   *Input
}

// NewMenu creates and initializes a Menu/GameState object
func NewMenu() *Menu {
	s := &Menu{input: NewInput()}

	s.widgets = []Widget{
		NewLabel("CAN YOU HERD KITTENS?", Acme.large),
		NewLabel("Move the yellow blob by clicking where you want it to go", Acme.normal),
		NewLabel("Build/demolish walls using the WASD keys", Acme.normal),
		NewLabel("Herd the kittens into the square in the middle", Acme.normal),
		NewTextButton("START", 200, 50, Acme.normal, func() { GSM.Switch(NewCutscene(7, 5, 4)) }),
		// NewTextButton("NORMAL", 200, 50, Acme.normal, func() { GSM.Switch(NewCutscene(15, 11, 4)) }),
		// NewTextButton("LARGE", 200, 50, Acme.large, func() { GSM.Switch(NewCutscene(21, 17, 8)) }),
		// NewTextButton("EXCESSIVE", 200, 50, Acme.normal, func() { GSM.Switch(NewCutscene(31, 23, 8)) }),
	}

	return s
}

// Layout implements ebiten.Game's Layout
func (s *Menu) Layout(outsideWidth, outsideHeight int) (int, int) {

	xCenter := outsideWidth / 2
	// create len(widgets) + 1 vertical slots
	yPlaces := [6]int{} // golang gotcha: can't use len(s.widgets)
	for i := 0; i < len(yPlaces); i++ {
		yPlaces[i] = (outsideHeight / len(yPlaces)) * i
	}

	for i, w := range s.widgets {
		w.SetPosition(xCenter, yPlaces[i+1])
	}

	return outsideWidth, outsideHeight
}

// Update updates the current game state.
func (s *Menu) Update() error {

	s.input.Update()

	if inpututil.IsKeyJustReleased(ebiten.KeyBackspace) {
		os.Exit(0)
	}

	for _, w := range s.widgets {
		if w.Pushed(s.input) {
			w.Action()
			break
		}
	}

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
