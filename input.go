package seGUI

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Input update, returns if it ate: left click, right click, or ebiten.key
func InputUpdate() (bool, bool, int) {
	mx, my := ebiten.CursorPosition()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {

		//Detect clicks within open windows
		for _, win := range openWindows {
			if posWithinRect(V2i{X: mx, Y: my}, win.bounds) {

				log.Printf("click within window: %v.\n", win.win.Title)
				return true, false, 0
			}
		}
	}
	return false, false, 0
}
