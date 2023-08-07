// Copyright ©️ 2021 oddstream.games

package maze

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gomaze/util"
)

var (
	ghostImages        map[int]*ebiten.Image
	directionlessImage *ebiten.Image
)

func init() {
	ghostImages = make(map[int]*ebiten.Image, 4)
	for d := 0; d < 4; d++ {
		img := makeGhostImage(d)
		ghostImages[d] = ebiten.NewImageFromImage(img)
	}
	directionlessImage = ebiten.NewImageFromImage(makeGhostImage(-1))
}

type dirfunc func(int) int

// Ghost defines the yellow blob/player avatar
type Ghost struct {
	tile                   *Tile   // tile we are sitting on
	dest                   *Tile   // tile we are lerping to
	facing                 int     // 0,1,2,3
	srcX, srcY, dstX, dstY float64 // positions for lerp
	lerpstep               float64
	speed                  float64
	worldX, worldY         float64
}

// NewGhost creates a new Ghost object
func NewGhost(start *Tile) *Ghost {
	gh := &Ghost{tile: start}
	gh.facing = rand.Intn(3)
	gh.speed = 0.01 + (rand.Float64() * 0.01)
	gh.worldX, gh.worldY = gh.tile.Position()
	return gh
}

// func (gh *Ghost) isPuckVisible(d int) bool {
// 	for t := gh.tile; !t.IsWall(d); t = t.Neighbour(d) {
// 		if t == TheGrid.puck.tile || t == TheGrid.puck.dest {
// 			return true
// 		}
// 	}
// 	return false
// }

func (gh *Ghost) isDirOkay(d int) bool {
	// can't go through walls
	if gh.tile.IsWall(d) {
		return false
	}

	tn := gh.tile.Neighbour(d)

	// don't like puck
	if tn == TheGrid.puck.tile || tn == TheGrid.puck.dest {
		return false
	}

	// don't leave the pen - this makes it too easy
	// if gh.tile.pen {
	// 	tn := gh.tile.Neighbour(dir)
	// 	if !tn.pen {
	// 		return false
	// 	}
	// }

	// don't like going on top of other ghosts
	// for _, g := range TheGrid.ghosts {
	// 	if g == gh {
	// 		continue
	// 	}
	// 	if g.tile == tn || g.dest == tn {
	// 		return false
	// 	}
	// }

	return true
}

// Update the state/position of the Ghost
func (gh *Ghost) Update() error {

	if gh.dest == nil {
		var dirfuncs [4]dirfunc
		if rand.Float64() < 0.5 {
			dirfuncs = [4]dirfunc{util.Leftward, util.Forward, util.Rightward, util.Backward}
		} else {
			dirfuncs = [4]dirfunc{util.Rightward, util.Forward, util.Leftward, util.Backward}
		}
		for d := 0; d < 4; d++ {
			newd := dirfuncs[d](gh.facing)
			if gh.isDirOkay(newd) {
				gh.facing = newd
				gh.dest = gh.tile.Neighbour(newd)
				break
			}
		}
		if gh.dest != nil {
			gh.lerpstep = 0
			gh.srcX, gh.srcY = gh.tile.Position()
			gh.dstX, gh.dstY = gh.dest.Position()
		}
		// if gh.dest == nil, ghost has no direction
		// this can happen when trying to stop ghosts from sitting on top of each other in Ghost.IsGoodDir()
	} else {
		if gh.lerpstep >= 1 {
			gh.tile = gh.dest
			gh.worldX, gh.worldY = gh.tile.Position()
			gh.dest = nil
		} else {
			gh.worldX = util.Lerp(gh.srcX, gh.dstX, gh.lerpstep)
			gh.worldY = util.Lerp(gh.srcY, gh.dstY, gh.lerpstep)
			gh.lerpstep += gh.speed
		}
	}

	return nil
}

// Draw the Ghost
func (gh *Ghost) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(gh.worldX, gh.worldY)
	op.GeoM.Translate(CameraX, CameraY)
	if gh.dest == nil {
		screen.DrawImage(directionlessImage, op)
	} else {
		screen.DrawImage(ghostImages[gh.facing], op)
	}
}
