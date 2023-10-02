package maze

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

// Label is a widget object that shows some text on screen
type Label struct {
	text          string
	font          font.Face
	origin        image.Point
	width, height int
}

// NewLabel creates and returns a new Label object
func NewLabel(str string, btnFont font.Face) *Label {
	l := &Label{text: str, font: btnFont}
	bound, _ := font.BoundString(l.font, l.text)
	l.width = (bound.Max.X - bound.Min.X).Ceil()
	l.height = (bound.Max.Y - bound.Min.Y).Ceil()
	return l
}

// SetPosition sets the position of this widget in screen coords
func (l *Label) SetPosition(x, y int) {
	bound, _ := font.BoundString(l.font, l.text)
	l.width = (bound.Max.X - bound.Min.X).Ceil()
	l.height = (bound.Max.Y - bound.Min.Y).Ceil()
	l.origin = image.Point{X: x - (l.width / 2), Y: y - (l.height / 2)}
}

// Rect gives the x,y coords of the label's top left and bottom right corners, in screen coordinates
func (l *Label) Rect() (x0 int, y0 int, x1 int, y1 int) {
	x0 = l.origin.X
	y0 = l.origin.Y
	x1 = x0 + l.width
	y1 = y0 + l.height
	return // using named return parameters
}

// Action invokes the action func
func (l *Label) Action() {
	// Labels take no action
}

// Update the button state (transitions, NOT user input)
func (l *Label) Update() error {
	return nil
}

// Draw handles rendering of Label object
func (l *Label) Draw(screen *ebiten.Image) {

	text.Draw(screen, l.text, l.font, l.origin.X, l.origin.Y+l.height, BasicColors["White"])

}
