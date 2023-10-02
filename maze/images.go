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

	dc.SetRGBA(0, 0, 0, 1) // black walls, get recolored when drawn
	dc.SetLineWidth(lineWidth)
	dc.SetLineCap(gg.LineCapRound)

	switch walls {
	case 0:
		// explicitly do nothing
	case NORTH_WALL:
		dc.DrawLine(w, n, e, n)
	case EAST_WALL:
		dc.DrawLine(e, n, e, s)
	case SOUTH_WALL:
		dc.DrawLine(w, s, e, s)
	case WEST_WALL:
		dc.DrawLine(w, n, w, s)

	case NORTH_WALL | SOUTH_WALL:
		dc.DrawLine(w, n, e, n)
		dc.DrawLine(w, s, e, s)
	case EAST_WALL | WEST_WALL:
		dc.DrawLine(e, n, e, s)
		dc.DrawLine(w, n, w, s)

	case NORTH_WALL | EAST_WALL:
		dc.MoveTo(w, n)
		dc.LineTo(e, n)
		dc.LineTo(e, s)
	case EAST_WALL | SOUTH_WALL:
		dc.MoveTo(e, n)
		dc.LineTo(e, s)
		dc.LineTo(w, s)
	case SOUTH_WALL | WEST_WALL:
		dc.MoveTo(e, s)
		dc.LineTo(w, s)
		dc.LineTo(w, n)
	case WEST_WALL | NORTH_WALL:
		dc.MoveTo(w, s)
		dc.LineTo(w, n)
		dc.LineTo(e, n)

	case NORTH_WALL | EAST_WALL | SOUTH_WALL:
		dc.MoveTo(w, n)
		dc.LineTo(e, n)
		dc.LineTo(e, s)
		dc.LineTo(w, s)
	case EAST_WALL | SOUTH_WALL | WEST_WALL:
		dc.MoveTo(e, n)
		dc.LineTo(e, s)
		dc.LineTo(w, s)
		dc.LineTo(w, n)
	case SOUTH_WALL | WEST_WALL | NORTH_WALL:
		dc.MoveTo(e, s)
		dc.LineTo(w, s)
		dc.LineTo(w, n)
		dc.LineTo(e, n)
	case WEST_WALL | NORTH_WALL | EAST_WALL:
		dc.MoveTo(w, s)
		dc.LineTo(w, n)
		dc.LineTo(e, n)
		dc.LineTo(e, s)

	case NORTH_WALL | EAST_WALL | SOUTH_WALL | WEST_WALL:
		dc.MoveTo(w, n)
		dc.LineTo(e, n)
		dc.LineTo(e, s)
		dc.LineTo(w, s)
		dc.ClosePath()

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

func makeNPCImage(dir int, bodyColor string) image.Image {
	mid := float64(TileSize / 2)
	dc := gg.NewContext(TileSize, TileSize)

	if col, ok := BasicColors[bodyColor]; ok {
		dc.SetColor(col)
		dc.MoveTo(polyCoords[0]*2+mid, polyCoords[1]*2+mid)
		for i := 2; i < len(polyCoords); {
			x := polyCoords[i]
			i++
			y := polyCoords[i]
			i++
			dc.LineTo(x*2+mid, y*2+mid)
		}
		dc.ClosePath()
		dc.Fill()
	}

	dc.SetColor(BasicColors["White"])
	dc.DrawCircle(mid-10, mid-2, 8)
	dc.DrawCircle(mid+10, mid-2, 8)
	dc.Fill()
	dc.SetColor(BasicColors["Black"])
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
