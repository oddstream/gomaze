package util

import (
	"image"
)

// InRect returns true if px,py is within Rect returned by function parameter
func InRect(pt image.Point, fn func() (int, int, int, int)) bool {
	x0, y0, x1, y1 := fn()
	return pt.X > x0 && pt.Y > y0 && pt.X < x1 && pt.Y < y1
}

// Lerp see https://en.wikipedia.org/wiki/Linear_interpolation
func Lerp(v0 float64, v1 float64, t float64) float64 {
	return (1-t)*v0 + t*v1
}

// Smoothstep see http://sol.gfxile.net/interpolation/
func Smoothstep(A float64, B float64, v float64) float64 {
	v = (v) * (v) * (3 - 2*(v)) // smoothstep
	// v = (v) * (v) * (v) * ((v)*((v)*6-15) + 10)	// smootherstep
	X := (B * v) + (A * (1.0 - v))
	return X
}

// Normalize is the opposite of lerp. Instead of a range and a factor, we give a range and a value to find out the factor.
func Normalize(start, finish, value float64) float64 {
	return (value - start) / (finish - start)
}

// MapValue converts a value from the scale [fromMin, fromMax] to a value from the scale [toMin, toMax].
// It’s just the normalize and lerp functions working together.
func MapValue(value, fromMin, fromMax, toMin, toMax float64) float64 {
	return Lerp(toMin, toMax, Normalize(fromMin, fromMax, value))
}

// Clamp a value between min and max values
// func Clamp(value, min, max float64) float64 {
// 	return math.Min(math.Max(value, min), max)
// }

// https://stackoverflow.com/questions/51626905/drawing-circles-with-two-radius-in-golang
// https://en.wikipedia.org/wiki/Midpoint_circle_algorithm
// func drawCircle(img *ebiten.Image, x0, y0, r int, c color.Color) {
// 	x, y, dx, dy := r-1, 0, 1, 1
// 	err := dx - (r * 2)

// 	for x > y {
// 		img.Set(x0+x, y0+y, c)
// 		img.Set(x0+y, y0+x, c)
// 		img.Set(x0-y, y0+x, c)
// 		img.Set(x0-x, y0+y, c)
// 		img.Set(x0-x, y0-y, c)
// 		img.Set(x0-y, y0-x, c)
// 		img.Set(x0+y, y0-x, c)
// 		img.Set(x0+x, y0-y, c)

// 		if err <= 0 {
// 			y++
// 			err += dy
// 			dy += 2
// 		}
// 		if err > 0 {
// 			x--
// 			dx += 2
// 			err += dx - (r * 2)
// 		}
// 	}
// }

func sign(p1, p2, p3 image.Point) int {
	return (p1.X-p3.X)*(p2.Y-p3.Y) - (p2.X-p3.X)*(p1.Y-p3.Y)
}

// https://stackoverflow.com/questions/2049582/how-to-determine-if-a-point-is-in-a-2d-triangle#2049593
func PointInTriangle(pt, v1, v2, v3 image.Point) bool {
	d1 := sign(pt, v1, v2)
	d2 := sign(pt, v2, v3)
	d3 := sign(pt, v3, v1)
	has_neg := (d1 < 0) || (d2 < 0) || (d3 < 0)
	has_pos := (d1 > 0) || (d2 > 0) || (d3 > 0)
	return !(has_neg && has_pos)
}
