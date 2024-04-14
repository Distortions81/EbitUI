package seGUI

import (
	"errors"
	"image"
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	viewerWidth, viewerHeight int
	windowList                map[string]*windowObject
	openWindows               []*windowObject

	windowsLock sync.Mutex

	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)

	mplusFaceSource *text.GoTextFaceSource
)

func getBounds(win *windowObject) FourV2i {
	rect := FourV2i{
		TopLeft:     V2i{X: win.position.X, Y: win.position.Y},
		TopRight:    V2i{X: win.position.X, Y: win.position.Y},
		BottomLeft:  V2i{X: win.position.X + win.size.X, Y: win.position.Y + win.size.Y},
		BottomRight: V2i{X: win.position.X + win.size.X, Y: win.position.Y + win.size.Y},
	}
	return rect
}

func updateWinPos(win *windowObject, pos V2i) {
	win.position = pos
	win.bounds = getBounds(win)
}

func updateWinSize(win *windowObject, size V2i) {
	win.size = size
	win.bounds = getBounds(win)
}

func posWithinRect(pos V2i, rect FourV2i) bool {
	if pos.X >= rect.TopLeft.X &&
		pos.Y >= rect.TopLeft.Y &&
		pos.X <= rect.BottomRight.X &&
		pos.Y <= rect.BottomRight.Y {
		return true
	}
	return false
}

func closeWindow(windowID string) error {
	windowID = strings.ToLower(windowID)

	window := windowList[windowID]

	if window != nil {
		if window.open {
			window.open = false

			numOpen := len(openWindows) - 1
			for w := numOpen; numOpen > 0; numOpen-- {
				if openWindows[w].id != windowID {

					//Delete item
					openWindows = append(openWindows[:w], openWindows[w+1:]...)
				}
			}
		}
		return nil
	}

	return errors.New("unable to find window")
}

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
