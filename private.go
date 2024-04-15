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
	windowList                map[WindowID]*windowObject
	openWindows               []*windowObject

	windowsLock sync.Mutex

	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)

	mplusFaceSource *text.GoTextFaceSource
)

func (posA V2i) addPos(posB V2i) V2i {
	return V2i{X: posA.X + posB.X, Y: posA.Y + posB.Y}
}

func (posA V2i) subPos(posB V2i) V2i {
	return V2i{X: posA.X - posB.X, Y: posA.Y - posB.Y}
}

func (win *windowObject) updateWin() {

	if win.position.X+win.size.X > viewerWidth {
		win.position.X = (viewerWidth - win.size.X)
	}
	if win.position.Y+win.size.Y > viewerHeight {
		win.position.Y = (viewerHeight - win.size.Y)
	}
	if win.position.X < 0 {
		win.position.X = 0
	}
	if win.position.Y < 0 {
		win.position.Y = 0
	}

	if win.size != win.oldSize {
		win.oldSize = win.size
		win.drawCache = ebiten.NewImage(win.size.X, win.size.Y)
		win.dirty = true
	}

	win.bounds = FourV2i{
		TopLeft:     V2i{X: win.position.X, Y: win.position.Y},
		TopRight:    V2i{X: win.size.X + win.position.X, Y: win.position.Y},
		BottomLeft:  V2i{X: win.position.X, Y: win.size.Y + win.position.Y},
		BottomRight: V2i{X: win.size.X + win.position.X, Y: win.size.Y + win.position.Y},
	}
	win.titleBounds = FourV2i{
		TopLeft:     win.bounds.TopLeft,
		TopRight:    win.bounds.TopRight,
		BottomLeft:  V2i{X: win.bounds.TopLeft.X, Y: win.bounds.TopLeft.Y + win.win.TitleSize},
		BottomRight: V2i{X: win.bounds.BottomRight.X, Y: win.bounds.TopLeft.Y + win.win.TitleSize},
	}
}

func (rect FourV2i) posWithinRect(pos V2i) bool {
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

	window := windowList[WindowID(windowID)]

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

func (win *windowObject) updateWindowCache() {
	if win.dirty {
		win.dirty = false
		win.drawCache.Fill(win.win.BGColor)
	}
}

// Call this in Ebiten Layout
func clampWindows() {
	windowsLock.Lock()
	defer windowsLock.Unlock()

	width, height := viewerWidth, viewerHeight

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

		win.updateWin()

		if changedSize {
			win.drawCache = ebiten.NewImage(win.size.X, win.size.Y)
			win.dirty = true
		}
	}
}
