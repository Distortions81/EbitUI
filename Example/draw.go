package main

import (
	"github.com/Distortions81/EbitUI"
	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) Draw(screen *ebiten.Image) {
	EbitUI.DrawWindows(screen)
}
