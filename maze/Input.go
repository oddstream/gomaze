package maze

import (
	"image"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type (
	// Observable https://gist.github.com/patrickmn/1549985
	Observable interface {
		Add(observer Observer)
		Notify(event interface{})
		Remove(event interface{})
	}

	// Observer https://gist.github.com/patrickmn/1549985
	Observer interface {
		NotifyCallback(event interface{})
	}
)

// Input records state of mouse and touch, Subject in Observer pattern
type Input struct {
	// pressed        map[ebiten.Key]struct{} // an empty and useless type
	observer sync.Map
}

// NewInput Input object constructor
func NewInput() *Input {
	// no fields to initialize, so use the built-in new()
	return new(Input)
}

// Add this observer to the list
func (i *Input) Add(observer Observer) {
	i.observer.Store(observer, struct{}{})
}

// Remove this observer from the list
func (i *Input) Remove(observer Observer) {
	i.observer.Delete(observer)
}

// Notify observers that an event has happened
func (i *Input) Notify(event interface{}) {
	i.observer.Range(func(key, value interface{}) bool {
		if key == nil {
			return false
		}
		key.(Observer).NotifyCallback(event)
		return true
	})
}

// Update the state of the Input object
func (i *Input) Update() {

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		i.Notify(image.Point{X: x, Y: y})
	} else {
		ts := inpututil.JustPressedTouchIDs()
		if len(ts) == 1 { // len() for nil slices is zero
			if inpututil.IsTouchJustReleased(ts[0]) {
				x, y := ebiten.TouchPosition(ts[0])
				i.Notify(image.Point{X: x, Y: y})
			}
		}
	}

	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if inpututil.IsKeyJustPressed(k) {
			i.Notify(k)
		}
	}
}
