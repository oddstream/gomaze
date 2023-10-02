package maze

import (
	"math/rand"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gomaze/util"
)

type Meanie struct {
	tile                   *Tile   // tile we are sitting on
	dest                   *Tile   // tile we are lerping to
	srcX, srcY, dstX, dstY float64 // positions for lerp
	lerpstep               float64
	speed                  float64
	worldX, worldY         float64
	img                    *ebiten.Image
}

func NewMeanie(start *Tile) *Meanie {
	m := &Meanie{tile: start}
	m.worldX, m.worldY = m.tile.position()
	m.speed = 0.01 + (rand.Float64() * 0.01)
	{
		dc := gg.NewContext(TileSize, TileSize)
		dc.SetColor(BasicColors["Black"])
		dc.DrawCircle(float64(TileSize/2), float64(TileSize/2), float64(TileSize/3))
		dc.Fill()
		dc.Stroke()
		m.img = ebiten.NewImageFromImage(dc.Image())
	}
	return m
}

func (m *Meanie) Update() error {
	if m.dest == nil {
		m.dest = TheGrid.findTileTowards(m.tile, TheGrid.puck.tile)
		if m.dest != nil {
			m.lerpstep = 0
			m.srcX, m.srcY = m.tile.position()
			m.dstX, m.dstY = m.dest.position()
		}
	} else {
		if m.lerpstep >= 1 {
			m.tile = m.dest
			m.worldX, m.worldY = m.tile.position()
			m.dest = nil
		} else {
			m.worldX = util.Lerp(m.srcX, m.dstX, m.lerpstep)
			m.worldY = util.Lerp(m.srcY, m.dstY, m.lerpstep)
			m.lerpstep += m.speed
		}
	}
	return nil
}

func (m *Meanie) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(m.worldX, m.worldY)
	op.GeoM.Translate(CameraX, CameraY)
	screen.DrawImage(m.img, op)
}
