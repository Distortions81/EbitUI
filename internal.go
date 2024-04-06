package seGUI

import "github.com/hajimehoshi/ebiten/v2"

func updateWindowCache(win *windowObject) {
	if win.dirty {
		win.dirty = false
		win.drawCache.Fill(win.win.BGColor)
	}
}

// Call this in Ebiten Layout
func clampWindows(width, height int) {
	windowsLock.Lock()
	defer windowsLock.Unlock()

	if width == viewerWidth && height == viewerHeight {
		return
	}

	viewerWidth, viewerHeight = width, height

	changedSize := false
	for w, win := range windowList {
		if win.id == "hud" {
			win.drawCache = ebiten.NewImage(width, height)
			win.dirty = true
			windowList[w].size = V2i{X: width, Y: height}
			continue
		}

		if win.win.Resizable {
			if win.size.X > width {
				win.size.X = width
				changedSize = true
			}
			if win.size.Y > height {
				win.size.Y = height
				changedSize = true
			}
		}

		if win.position.X+win.size.X > width {
			win.position.X = (width - win.size.X)
		}
		if win.position.Y+win.size.Y > height {
			win.position.Y = (height - win.size.Y)
		}

		if changedSize {
			win.drawCache = ebiten.NewImage(win.size.X, win.size.Y)
			win.dirty = true
		}
	}
}
