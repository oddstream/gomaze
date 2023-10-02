package maze

import (
	"fmt"
	"image/color"
	"log"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gomaze/util"
)

// Puck defines the yellow blob/player avatar
type Puck struct {
	tile                   *Tile   // tile we are sitting on
	dest                   *Tile   // tile we are lerping to
	srcX, srcY, dstX, dstY float64 // positions for lerp
	lerpstep               float64
	img                    *ebiten.Image
	worldX, worldY         float64
	col                    color.RGBA
	ball                   *Ball
}

// NewPuck creates a new Puck object
func NewPuck(start *Tile, col color.RGBA) *Puck {
	p := &Puck{tile: start, col: col}
	p.createImage()
	p.worldX, p.worldY = p.tile.position()
	p.SetCamera()
	p.ball = NewBall(start, col)
	return p
}

func (p *Puck) createImage() {
	dc := gg.NewContext(TileSize, TileSize)
	dc.SetColor(p.col)
	dc.DrawCircle(float64(TileSize/2), float64(TileSize/2), float64(TileSize/3))
	dc.Fill()
	dc.Stroke()
	p.img = ebiten.NewImageFromImage(dc.Image())
}

// SetCamera so that puck is at the center of the screen
func (p *Puck) SetCamera() {
	w, h := WindowWidth, WindowHeight // can't use ebiten.WindowSize(), returns 0,0 on WASM
	// sx, sy := p.puckImage.Size()
	sx := p.img.Bounds().Dx()
	sy := p.img.Bounds().Dy()
	CameraX = float64(w/2-sx/2) - p.worldX
	CameraY = float64(h/2-sy/2) - p.worldY
}

// ThrowBallTo a target tile
func (p *Puck) ThrowBallTo(targ *Tile) {

	// if puck is lerping, stop it
	if p.dest != nil {
		p.tile = p.dest
		p.worldX, p.worldY = p.tile.position()
		p.dest = nil
	}

	TheGrid.BFS(p.tile, targ)
	for t := targ; t != nil && t != p.tile; t = t.parent {
		t.marked = true
	}
	p.ball.StartThrow(targ)
}

// String representation of puck
func (p *Puck) String() string {
	return fmt.Sprintf("(%v,%v)", p.tile.X, p.tile.Y)
}

// Update the state/position of the Puck and it's ball
func (p *Puck) Update() error {

	p.ball.Update()

	if p.tile.marked {
		// println("unmarking")
		p.tile.marked = false
	}

	if p.dest == nil {
		// if any of the neighbours are marked, move there
		for _, d := range ALL_DIRECTIONS {
			if p.tile.isWall(d) {
				continue
			}
			tn := p.tile.neighbour(d)
			if tn == nil {
				log.Fatal("unwalled edge found in Puck Update")
			}
			if tn.marked {
				p.dest = tn
				p.srcX, p.srcY = p.tile.position()
				p.dstX, p.dstY = p.dest.position()
				p.lerpstep = 0
				break
			}
		}
	} else {
		if p.lerpstep >= 1 {
			p.tile = p.dest
			p.tile.marked = false
			p.worldX, p.worldY = p.tile.position()
			p.dest = nil
		} else {
			p.worldX = util.Lerp(p.srcX, p.dstX, p.lerpstep)
			p.worldY = util.Lerp(p.srcY, p.dstY, p.lerpstep)
			p.lerpstep += 0.05
		}
	}
	p.SetCamera()

	return nil
}

// Draw the Puck, and maybe it's ball
func (p *Puck) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.worldX, p.worldY)
	op.GeoM.Translate(CameraX, CameraY)
	screen.DrawImage(p.img, op)
	if p.tile != p.ball.tile {
		p.ball.Draw(screen)
	}
}
