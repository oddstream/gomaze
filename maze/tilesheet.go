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
