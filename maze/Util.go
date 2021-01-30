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

// SMOOTHSTEP from http://sol.gfxile.net/interpolation/
func SMOOTHSTEP(x float64) float64 {
	return (x) * (x) * (3 - 2*(x))
}

// SMOOTHERSTEP from http://sol.gfxile.net/interpolation/
func SMOOTHERSTEP(x float64) float64 {
	return (x) * (x) * (x) * ((x)*((x)*6-15) + 10)
}

func smoothstep(A float64, B float64, v float64) float64 {
	v = SMOOTHSTEP(v)
	X := (B * v) + (A * (1.0 - v))
	return X
}

/*
// https://en.wikipedia.org/wiki/Smoothstep
func smoothstep(edge0 float64, edge1 float64, x float64) float64 {
	// Scale, bias and saturate x to 0..1 range
	x = clamp((x-edge0)/(edge1-edge0), 0.0, 1.0)
	// Evaluate polynomial
	return x * x * (3 - 2*x)
}

func smootherstep(edge0 float64, edge1 float64, x float64) float64 {
	// Scale, and clamp x to 0..1 range
	x = clamp((x-edge0)/(edge1-edge0), 0.0, 1.0)
	// Evaluate polynomial
	return x * x * x * (x*(x*6-15) + 10)
}

func clamp(x float64, lowerlimit float64, upperlimit float64) float64 {
	if x < lowerlimit {
		x = lowerlimit
	}
	if x > upperlimit {
		x = upperlimit
	}
	return x
}
*/

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

func opposite(dir int) int {
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
