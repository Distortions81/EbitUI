package seGUI

import (
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
			if win.titleBounds.posWithinRect(mousePos) {

				mouseDiff := mousePos.subPos(lastMouse)
				win.position = V2i{X: win.position.X + mouseDiff.X, Y: win.position.Y + mouseDiff.Y}
				win.updateWin()
				return true, false, 0
			}
		}

	}

	return false, false, 0
}
