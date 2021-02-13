// Copyright ©️ 2021 oddstream.games

package maze

import (
	"image"
	"log"

	"github.com/fogleman/gg"
)

// return an image.Image that is bigger than the tile size requested so endcaps are visible
func makeTileImage(walls uint, unreachable bool) image.Image {

	tileSizeEx := TileSize + (TileSize / 6) // same as linewidth

	margin := float64(tileSizeEx-TileSize) / 2
	lineWidth := float64(TileSize / 6)

	n := margin
	e := float64(tileSizeEx) - margin
	s := float64(tileSizeEx) - margin
	w := margin

	dc := gg.NewContext(tileSizeEx, tileSizeEx)

	if !unreachable {
		dc.SetRGBA(0, 0, 0, 0.2)
		dc.DrawRectangle(margin, margin, float64(TileSize), float64(TileSize))
		dc.Fill()
		dc.Stroke()
	}

	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(lineWidth)
	dc.SetLineCap(gg.LineCapRound)

	switch walls {
	case 0:
		// explicitly do nothing
	case NORTH:
		dc.DrawLine(w, n, e, n)
	case EAST:
		dc.DrawLine(e, n, e, s)
	case SOUTH:
		dc.DrawLine(w, s, e, s)
	case WEST:
		dc.DrawLine(w, n, w, s)

	case NORTH | SOUTH:
		dc.DrawLine(w, n, e, n)
		dc.DrawLine(w, s, e, s)
	case EAST | WEST:
		dc.DrawLine(e, n, e, s)
		dc.DrawLine(w, n, w, s)

	case NORTH | EAST:
		dc.MoveTo(w, n)
		dc.LineTo(e, n)
		dc.LineTo(e, s)
	case EAST | SOUTH:
		dc.MoveTo(e, n)
		dc.LineTo(e, s)
		dc.LineTo(w, s)
	case SOUTH | WEST:
		dc.MoveTo(e, s)
		dc.LineTo(w, s)
		dc.LineTo(w, n)
	case WEST | NORTH:
		dc.MoveTo(w, s)
		dc.LineTo(w, n)
		dc.LineTo(e, n)

	case NORTH | EAST | SOUTH:
		dc.MoveTo(w, n)
		dc.LineTo(e, n)
		dc.LineTo(e, s)
		dc.LineTo(w, s)
	case EAST | SOUTH | WEST:
		dc.MoveTo(e, n)
		dc.LineTo(e, s)
		dc.LineTo(w, s)
		dc.LineTo(w, n)
	case SOUTH | WEST | NORTH:
		dc.MoveTo(e, s)
		dc.LineTo(w, s)
		dc.LineTo(w, n)
		dc.LineTo(e, n)
	case WEST | NORTH | EAST:
		dc.MoveTo(w, s)
		dc.LineTo(w, n)
		dc.LineTo(e, n)
		dc.LineTo(e, s)

	case NORTH | EAST | SOUTH | WEST:
		dc.MoveTo(w, n)
		dc.LineTo(e, n)
		dc.LineTo(e, s)
		dc.LineTo(w, s)
		dc.ClosePath()
		// dc.FillPreserve()

		// dc.DrawRoundedRectangle(w, n, float64(tileSize), float64(tileSize), lineWidth)
	default:
		log.Fatal("makeTile called with wrong bits", walls)
	}
	dc.Stroke()

	return dc.Image()
}

var polyCoords = []float64{
	-12, 10, // bottom left
	-8, 8,
	-4, 12,
	0, 8,
	4, 12,
	8, 8,
	12, 10, // bottom right
	12, 0,
	// arch
	11, -5,
	10, -7,
	9, -8,
	8, -9,
	6, -10,

	0, -11,

	-6, -10,
	-8, -9,
	-9, -8,
	-10, -7,
	-11, -5,
	// end of arch
	-12, 0,
}

func makeGhostImage(dir int) image.Image {
	mid := float64(TileSize / 2)
	dc := gg.NewContext(TileSize, TileSize)

	// comment this out to just draw googly eyes, which is kinda fun
	// dc.SetColor(BasicColors["Silver"])
	// dc.MoveTo(polyCoords[0]*2+mid, polyCoords[1]*2+mid)
	// for i := 2; i < len(polyCoords); {
	// 	x := polyCoords[i]
	// 	i++
	// 	y := polyCoords[i]
	// 	i++
	// 	dc.LineTo(x*2+mid, y*2+mid)
	// }
	// dc.ClosePath()
	// dc.Fill()

	dc.SetRGB(1, 1, 1)
	dc.DrawCircle(mid-10, mid-2, 8)
	dc.DrawCircle(mid+10, mid-2, 8)
	dc.Fill()
	dc.SetRGB(0, 0, 0)
	switch dir {
	case -1: // kludge for directionless
		dc.DrawCircle(mid-10, mid-2, 4)
		dc.DrawCircle(mid+10, mid-2, 4)
	case 0: // NORTH
		dc.DrawCircle(mid-10, mid-6, 4)
		dc.DrawCircle(mid+10, mid-6, 4)
	case 1: // EAST
		dc.DrawCircle(mid-10+4, mid-2, 4)
		dc.DrawCircle(mid+10+4, mid-2, 4)
	case 2: //SOUTH
		dc.DrawCircle(mid-10, mid+2, 4)
		dc.DrawCircle(mid+10, mid+2, 4)
	case 3: // WEST
		dc.DrawCircle(mid-10-4, mid-2, 4)
		dc.DrawCircle(mid+10-4, mid-2, 4)
	}
	dc.Fill()
	dc.Stroke()
	return dc.Image()
}
