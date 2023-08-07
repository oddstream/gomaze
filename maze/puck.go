// Copyright ©️ 2021 oddstream.games

package maze

import (
	"fmt"
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

	puckImage *ebiten.Image

	worldX, worldY float64

	ball *Ball
}

// NewPuck creates a new Puck object
func NewPuck(start *Tile) *Puck {
	p := &Puck{tile: start}

	dc := gg.NewContext(TileSize, TileSize)
	dc.SetColor(BasicColors["Yellow"])
	dc.DrawCircle(float64(TileSize/2), float64(TileSize/2), float64(TileSize/3))
	dc.Fill()
	dc.Stroke()
	p.puckImage = ebiten.NewImageFromImage(dc.Image())

	p.worldX, p.worldY = p.tile.Position()
	p.SetCamera()

	p.ball = NewBall(start)

	return p
}

// SetCamera so that puck is at the center of the screen
func (p *Puck) SetCamera() {
	w, h := WindowWidth, WindowHeight // can't use ebiten.WindowSize(), returns 0,0 on WASM
	// sx, sy := p.puckImage.Size()
	sx := p.puckImage.Bounds().Dx()
	sy := p.puckImage.Bounds().Dy()
	CameraX = float64(w/2-sx/2) - p.worldX
	CameraY = float64(h/2-sy/2) - p.worldY
}

// ThrowBallTo a target tile
func (p *Puck) ThrowBallTo(targ *Tile) {

	// if puck is lerping, stop it
	if p.dest != nil {
		p.tile = p.dest
		p.worldX, p.worldY = p.tile.Position()
		p.dest = nil
	}

	found := false
	q := []*Tile{p.tile}
	p.tile.parent = p.tile
	for len(q) > 0 {
		t := q[0]
		q = q[1:] // take first tile off front of queue
		// println("pop", t.X, t.Y, "len now", len(q))
		if t == targ {
			found = true
			for path := t; path != p.tile; path = path.parent {
				path.marked = true
			}
			break
		}
		for _, d := range []int{0, 1, 2, 3} {
			if t.IsWall(d) {
				continue
			}
			tn := t.Neighbour(d)
			if tn == nil {
				log.Fatal("open unwalled edge found in Puck BFS")
			}
			if tn.parent == nil {
				tn.parent = t
				q = append(q, tn)
				// println("push", tn.X, tn.Y, "len now", len(q))
			}
		}
	}

	if found {
		p.ball.StartThrow(targ)
	}
}

// String representation of puck
func (p *Puck) String() string {
	return fmt.Sprintf("(%v,%v)", p.tile.X, p.tile.Y)
}

// Update the state/position of the Puck
func (p *Puck) Update() error {

	p.ball.Update()

	if p.tile.marked {
		// println("unmarking")
		p.tile.marked = false
	}

	if p.dest == nil {
		// if any of the neighbours are marked, move there
		for d := 0; d < 4; d++ {
			if p.tile.IsWall(d) {
				continue
			}
			tn := p.tile.Neighbour(d)
			if tn == nil {
				log.Fatal("unwalled edge found in Puck Update")
			}
			if tn.marked {
				p.dest = tn
				p.srcX, p.srcY = p.tile.Position()
				p.dstX, p.dstY = p.dest.Position()
				p.lerpstep = 0
				break
			}
		}
	} else {
		if p.lerpstep >= 1 {
			p.tile = p.dest
			p.tile.marked = false
			p.worldX, p.worldY = p.tile.Position()
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

// Draw the Puck
func (p *Puck) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.worldX, p.worldY)
	op.GeoM.Translate(CameraX, CameraY)
	screen.DrawImage(p.puckImage, op)
	p.ball.Draw(screen)
}
