package maze

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"oddstream.games/gomaze/util"
)

var (
	ghostImages             map[int]*ebiten.Image
	directionlessGhostImage *ebiten.Image
)

func init() {
	ghostImages = make(map[int]*ebiten.Image, 4)
	for d := 0; d < 4; d++ {
		img := makeNPCImage(d, "")
		ghostImages[d] = ebiten.NewImageFromImage(img)
	}
	directionlessGhostImage = ebiten.NewImageFromImage(makeNPCImage(-1, ""))
}

// Forward returns the direction (0-3)
func Forward(dir Direction) Direction {
	return dir
}

// Backward returns the direction (0-3)
func Backward(dir Direction) Direction {
	d := [4]Direction{SOUTH, WEST, NORTH, EAST}
	return d[dir]
}

// Leftward returns the direction (0-3)
func Leftward(dir Direction) Direction {
	d := [4]Direction{WEST, NORTH, EAST, SOUTH}
	return d[dir]
}

// Rightward returns the direction (0-3)
func Rightward(dir Direction) Direction {
	d := [4]Direction{EAST, SOUTH, WEST, NORTH}
	return d[dir]
}

type dirfunc func(Direction) Direction

// Ghost is a direction-following NPC that tries to avoid the Puck
type Ghost struct {
	tile                   *Tile     // tile we are sitting on
	dest                   *Tile     // tile we are lerping to
	facing                 Direction // NORTH || EAST || SOUTH || WEST
	srcX, srcY, dstX, dstY float64   // positions for lerp
	lerpstep               float64
	speed                  float64
	worldX, worldY         float64
}

// NewGhost creates a new Ghost object
func NewGhost(start *Tile) *Ghost {
	gh := &Ghost{tile: start}
	gh.facing = Direction(rand.Intn(3))
	gh.speed = 0.01 + (rand.Float64() * 0.02)
	gh.worldX, gh.worldY = gh.tile.position()
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

func (gh *Ghost) isDirOkay(d Direction) bool {
	// can't go through walls
	if gh.tile.isWall(d) {
		return false
	}

	tn := gh.tile.neighbour(d)

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
			dirfuncs = [4]dirfunc{Leftward, Forward, Rightward, Backward}
		} else {
			dirfuncs = [4]dirfunc{Rightward, Forward, Leftward, Backward}
		}
		for d := 0; d < 4; d++ {
			newd := dirfuncs[d](gh.facing)
			if gh.isDirOkay(newd) {
				gh.facing = newd
				gh.dest = gh.tile.neighbour(newd)
				break
			}
		}
		if gh.dest != nil {
			gh.lerpstep = 0
			gh.srcX, gh.srcY = gh.tile.position()
			gh.dstX, gh.dstY = gh.dest.position()
		}
		// if gh.dest == nil, ghost has no direction
		// this can happen when trying to stop ghosts from sitting on top of each other in Ghost.IsGoodDir()
	} else {
		if gh.lerpstep >= 1 {
			gh.tile = gh.dest
			gh.worldX, gh.worldY = gh.tile.position()
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
		screen.DrawImage(directionlessGhostImage, op)
	} else {
		screen.DrawImage(ghostImages[int(gh.facing)], op)
	}
}
