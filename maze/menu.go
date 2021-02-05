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
		NewLabel("MAZE", Acme.large),
		NewTextButton("TINY", 200, 50, Acme.normal, func() { TheGrid = NewGrid(5, 5, 4); GSM.Switch(TheGrid) }),
		NewTextButton("NORMAL", 200, 50, Acme.normal, func() { TheGrid = NewGrid(13, 11, 4); GSM.Switch(TheGrid) }),
		NewTextButton("LARGE", 200, 50, Acme.large, func() { TheGrid = NewGrid(21, 21, 12); GSM.Switch(TheGrid) }),
		NewTextButton("EXCESSIVE", 200, 50, Acme.normal, func() { TheGrid = NewGrid(31, 31, 24); GSM.Switch(TheGrid) }),
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
