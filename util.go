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
	halfx, halfy := win.size.X/2, win.size.Y/2
	rect := FourV2i{
		TopLeft:     V2i{X: win.position.X - halfx, Y: win.position.Y - halfy},
		TopRight:    V2i{X: win.position.X + halfx, Y: win.position.Y - halfy},
		BottomLeft:  V2i{X: win.position.X - halfx, Y: win.position.Y + halfy},
		BottomRight: V2i{X: win.position.X + halfx, Y: win.position.Y + halfy},
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
