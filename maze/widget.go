// Copyright ©️ 2021 oddstream.games

package maze

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Widget is an interface for widget objects
type Widget interface {
	Update() error
	Draw(*ebiten.Image)
	SetPosition(int, int)
	Rect() (int, int, int, int)
	Action()
}
