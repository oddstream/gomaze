// Copyright ©️ 2020 oddstream.games

// go mod init oddstream.games/tetra
// go mod tidy

// the package defining a command (an excutable Go program) always has the name main
// this is a signal to go build that it must invoke the linker to make an executable file
package main

import (
	"flag"
	"log"
	"os"

	// load png decoder in main package
	// _ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	tetra "oddstream.games/tetra/tetra"
)

func init() {
	flag.BoolVar(&tetra.DebugMode, "debug", false, "turn debug graphics on")
	flag.IntVar(&tetra.WindowWidth, "width", 1920/2, "width of window in pixels")
	flag.IntVar(&tetra.WindowHeight, "height", 1080/2, "height of window in pixels")
}

func main() {
	flag.Parse()

	if tetra.DebugMode {
		for i, a := range os.Args {
			println(i, a)
		}
	}

	game, err := tetra.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowTitle("Tetra Loops")                        // does nothing when runtime.GOARCH == "wasm"
	ebiten.SetWindowSize(tetra.WindowWidth, tetra.WindowHeight) // does nothing when runtime.GOARCH == "wasm"
	ebiten.SetWindowResizable(true)                             // does nothing when runtime.GOARCH == "wasm"
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
