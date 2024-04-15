package seGUI

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var lastMouse V2i
var dragWindow string

// Input update, returns if it ate: left click, right click, or ebiten.key
func InputUpdate() (bool, bool, int) {

	if !ebiten.IsFocused() {
		return false, false, 0
	}

	mx, my := ebiten.CursorPosition()
	mousePos := V2i{X: mx, Y: my}
	mousePos.clampToViewer()

	defer func() {
		lastMouse = mousePos
	}()

	if inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft) > 1 {

		//Detect clicks within open windows
		for _, win := range openWindows {
			if dragWindow == win.id ||
				win.titleBounds.contains(mousePos) {
				mouseDiff := mousePos.subPos(lastMouse)

				dragWindow = win.id

				//No change, don't recalc
				if mouseDiff.X == 0 && mouseDiff.Y == 0 {
					return true, false, 0
				}

				win.position = V2i{X: win.position.X + mouseDiff.X, Y: win.position.Y + mouseDiff.Y}
				win.updateWin()
				return true, false, 0
			}
		}
	} else {
		dragWindow = ""
	}

	return false, false, 0
}
