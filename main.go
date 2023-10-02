// go mod init oddstream.games/gomaze
// go mod tidy

// the package defining a command (an excutable Go program) always has the name main
// this is a signal to go build that it must invoke the linker to make an executable file
package main

import (
	"flag"
	"log"
	"os"

	// load png decoder in main package
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	maze "oddstream.games/gomaze/maze"
)

func init() {
	flag.BoolVar(&maze.DebugMode, "debug", false, "turn debug graphics on")
	flag.IntVar(&maze.WindowWidth, "width", 1920/2, "width of window in pixels")
	flag.IntVar(&maze.WindowHeight, "height", 1080/2, "height of window in pixels")
}

func main() {
	flag.Parse()

	if maze.DebugMode {
		for i, a := range os.Args {
			println(i, a)
		}
	}

	game, err := maze.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowTitle("Maze")                             // does nothing when runtime.GOARCH == "wasm"
	ebiten.SetWindowSize(maze.WindowWidth, maze.WindowHeight) // does nothing when runtime.GOARCH == "wasm"
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
