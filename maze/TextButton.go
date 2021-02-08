// Copyright ©️ 2021 oddstream.games

package maze

import (
	"image"
	"os"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"oddstream.games/gomaze/util"
)

// TextButton is an object that represents a button, Observer in Observer pattern
type TextButton struct {
	text          string
	font          font.Face
	action        func()
	origin        image.Point
	width, height int
	img           *ebiten.Image
}

// NewTextButton creates and returns a new TextButton object centered at x,y
func NewTextButton(str string, w int, h int, btnFont font.Face, actionFn func(), i *Input) *TextButton {

	tb := &TextButton{text: str, width: w, height: h, font: btnFont, action: actionFn}

	dc := gg.NewContext(w, h)
	dc.SetRGB(0, 0, 0)
	dc.DrawRoundedRectangle(0, 0, float64(w), float64(h), float64(w/20))
	dc.Fill()
	dc.SetRGB(1, 1, 1)
	dc.SetFontFace(tb.font)
	dc.DrawStringAnchored(tb.text, float64(w/2), float64(h/2), 0.5, 0.333)
	dc.Stroke()
	tb.img = ebiten.NewImageFromImage(dc.Image())

	i.Add(tb)

	return tb
}

// NotifyCallback is called by the Subject (Input) when something interesting happens
func (tb *TextButton) NotifyCallback(event interface{}) {
	switch v := event.(type) { // Type switch https://tour.golang.org/methods/16
	case image.Point:
		pt := v
		if util.InRect(pt, tb.Rect) {
			tb.Action()
		}
	case ebiten.Key:
		k := v
		if k == ebiten.KeyBackspace {
			os.Exit(0)
		}
	}
}

// SetPosition sets the position of this widget in screen coords
func (tb *TextButton) SetPosition(x, y int) {
	tb.origin = image.Point{X: x - (tb.width / 2), Y: y - (tb.height / 2)}
}

// Rect gives the x,y coords of the TextButton's top left and bottom right corners, in screen coordinates
func (tb *TextButton) Rect() (x0 int, y0 int, x1 int, y1 int) {
	x0 = tb.origin.X
	y0 = tb.origin.Y
	x1 = x0 + tb.width
	y1 = y0 + tb.height
	return // using named return parameters
}

// Action invokes the action func
func (tb *TextButton) Action() {
	if tb.action != nil {
		tb.action()
	}
}

// Update the button state (transitions, NOT user input)
func (tb *TextButton) Update() error {
	return nil
}

// Draw handles rendering of TextButton object
func (tb *TextButton) Draw(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(tb.origin.X), float64(tb.origin.Y))
	screen.DrawImage(tb.img, op)

}
