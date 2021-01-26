// Copyright ©️ 2021 oddstream.games

package maze

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// InRect returns true of px,py is within Rect returned by function parameter
func InRect(pt image.Point, fn func() (int, int, int, int)) bool {
	x0, y0, x1, y1 := fn()
	return pt.X > x0 && pt.Y > y0 && pt.X < x1 && pt.Y < y1
}

// https://en.wikipedia.org/wiki/Linear_interpolation
func lerp(v0 float64, v1 float64, t float64) float64 {
	return (1-t)*v0 + t*v1
}

// https://stackoverflow.com/questions/51626905/drawing-circles-with-two-radius-in-golang
// https://en.wikipedia.org/wiki/Midpoint_circle_algorithm
func drawCircle(img *ebiten.Image, x0, y0, r int, c color.Color) {
	x, y, dx, dy := r-1, 0, 1, 1
	err := dx - (r * 2)

	for x > y {
		img.Set(x0+x, y0+y, c)
		img.Set(x0+y, y0+x, c)
		img.Set(x0-y, y0+x, c)
		img.Set(x0-x, y0+y, c)
		img.Set(x0-x, y0-y, c)
		img.Set(x0-y, y0-x, c)
		img.Set(x0+y, y0-x, c)
		img.Set(x0+x, y0-y, c)

		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - (r * 2)
		}
	}
}
