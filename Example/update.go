package main

import "github.com/Distortions81/EbitUI"

func (g *Game) Update() error {
	EbitUI.InputUpdate()
	return nil
}
