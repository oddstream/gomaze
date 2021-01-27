// Copyright ©️ 2021 oddstream.games

package maze

import (
	"log"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// StatePuckSettled when the Puck is not moving
	StatePuckSettled = iota
	// StatePuckMoving when Puck is lerping/chasing it's ball
	StatePuckMoving
)

// PuckState of this Puck
type PuckState int

// Puck defines the yellow blob/player avatar
type Puck struct {
	state PuckState

	tile                   *Tile   // tile we are sitting on
	home                   *Tile   // tile we started on
	dest                   *Tile   // tile we are lerping to
	srcX, srcY, dstX, dstY float64 // positions for lerp
	lerpstep               float64

	puckImage *ebiten.Image

	x, y float64

	ball *Ball
}

// NewPuck creates a new Puck object
func NewPuck(start *Tile) *Puck {
	p := &Puck{home: start, tile: start, state: StatePuckSettled}

	dc := gg.NewContext(TileSize, TileSize)
	dc.SetRGB(1, 1, 0) // Yellow
	dc.DrawCircle(float64(TileSize/2), float64(TileSize/2), float64(TileSize/3))
	dc.Fill()
	dc.Stroke()
	p.puckImage = ebiten.NewImageFromImage(dc.Image())

	p.x, p.y = p.tile.Position()

	p.ball = NewBall(start)

	return p
}

// ThrowBallTo a target tile
func (p *Puck) ThrowBallTo(targ *Tile) {
	p.ball.ThrowTo(targ)

	// if puck is lerping, stop it
	if p.dest != nil {
		p.tile = p.dest
		p.x, p.y = p.tile.Position()
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
		for d := 0; d < 4; d++ {
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
	if !found {
		log.Fatal("tile not found")
	}
	println("tile found")
}

// BallTile getter for location of puck's ball
// func (p *Puck) BallTile() *Tile {
// 	return p.ball.Tile()
// }

// Update the state/position of the Puck
func (p *Puck) Update() error {

	p.ball.Update()

	if p.tile.marked {
		println("unmarking")
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
				p.lerpstep = 0.1
				break
			}
		}
	} else {
		if p.lerpstep >= 1 {
			p.tile = p.dest
			p.tile.marked = false
			p.x, p.y = p.tile.Position()
			p.dest = nil
		} else {
			p.x = lerp(p.srcX, p.dstX, p.lerpstep)
			p.y = lerp(p.srcY, p.dstY, p.lerpstep)
			p.lerpstep += 0.1
		}
	}

	return nil
}

// Draw the Puck
func (p *Puck) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.x, p.y)
	screen.DrawImage(p.puckImage, op)
	p.ball.Draw(screen)
}
