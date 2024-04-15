package main

import (
	"os"
	"syscall"

	"github.com/Distortions81/EbitUI"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	signalHandle chan os.Signal
)

type Game struct {
}

func startEbiten() {

	// Set up ebiten
	ebiten.SetVsyncEnabled(true)
	ebiten.SetTPS(ebiten.SyncWithFPS)

	/* Set up our window */
	ebiten.SetWindowSize(defaultWindowWidth, defaultWindowHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("EbitUI Example")

	if err := ebiten.RunGameWithOptions(newGame(), nil); err != nil {
		return
	}

	signalHandle <- syscall.SIGINT
}

func newGame() *Game {

	return &Game{}
}

/* Window size chaged, handle it */
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return EbitUI.UpdateViewerSize(outsideWidth, outsideHeight)
}
