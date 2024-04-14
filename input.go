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
				titleRect := FourV2i{
					TopLeft:     win.bounds.TopLeft,
					TopRight:    win.bounds.TopRight,
					BottomLeft:  V2i{X: win.bounds.TopLeft.X, Y: win.bounds.TopLeft.Y + win.win.TitleSize},
					BottomRight: V2i{X: win.bounds.BottomRight.X, Y: win.bounds.TopLeft.Y + win.win.TitleSize},
				}
				if posWithinRect(mousePos, titleRect) {
					log.Printf("Window %v: Drag on titlebar.\n", win.win.Title)
					return true, false, 0
				}
			}
		}

	}

	return false, false, 0
}
