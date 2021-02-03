// Copyright ©️ 2021 oddstream.games

package maze

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// InRect returns true if px,py is within Rect returned by function parameter
func InRect(pt image.Point, fn func() (int, int, int, int)) bool {
	x0, y0, x1, y1 := fn()
	return pt.X > x0 && pt.Y > y0 && pt.X < x1 && pt.Y < y1
}

// https://en.wikipedia.org/wiki/Linear_interpolation
func lerp(v0 float64, v1 float64, t float64) float64 {
	return (1-t)*v0 + t*v1
}

func smoothstep(A float64, B float64, v float64) float64 {
	// http://sol.gfxile.net/interpolation/
	v = (v) * (v) * (3 - 2*(v)) // smoothstep
	// v = (v) * (v) * (v) * ((v)*((v)*6-15) + 10)	// smootherstep
	X := (B * v) + (A * (1.0 - v))
	return X
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

func forward(dir int) int {
	return dir
}

func backward(dir int) int {
	d := [4]int{2, 3, 0, 1}
	return d[dir]
}

func left(dir int) int {
	d := [4]int{3, 0, 1, 2}
	return d[dir]
}

func right(dir int) int {
	d := [4]int{1, 2, 3, 0}
	return d[dir]
}
