package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Tetris")

	app := newApp()

	if err := ebiten.RunGame(app); err != nil {
		log.Fatal(err)
	}
}
