package seGUI

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var lastMouse V2i

// Input update, returns if it ate: left click, right click, or ebiten.key
func InputUpdate() (bool, bool, int) {
	mx, my := ebiten.CursorPosition()
	mousePos := V2i{X: mx, Y: my}
	defer func() {
		lastMouse = mousePos
	}()

	if inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft) > 0 {

		//Detect clicks within open windows
		for _, win := range openWindows {
			if posWithinRect(mousePos, win.bounds) {
				if posWithinRect(mousePos, win.titleBounds) {
					log.Printf("Window %v: Drag on titlebar.\n", win.win.Title)

					mouseDiff := posDiff(mousePos, lastMouse)
					win.position = V2i{X: win.position.X + mouseDiff.X, Y: win.position.Y + mouseDiff.Y}
					break
				}
			}
		}

		return true, false, 0
	}

	return false, false, 0
}

func posDiff(posA, posB V2i) V2i {
	return V2i{X: posA.X - posB.X, Y: posA.Y - posB.Y}
}
