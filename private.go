package EbitUI

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

func (posA *V2i) addPos(posB V2i) V2i {
	return V2i{X: posA.X + posB.X, Y: posA.Y + posB.Y}
}

func (posA V2i) subPos(posB V2i) V2i {
	return V2i{X: posA.X - posB.X, Y: posA.Y - posB.Y}
}

func (pos *V2i) clampToViewer() {
	if pos.X < 0 {
		pos.X = 0
	} else if pos.X > viewerWidth {
		pos.X = viewerWidth
	}
	if pos.Y < 0 {
		pos.Y = 0
	} else if pos.Y > viewerHeight {
		pos.Y = viewerHeight
	}
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
		win.clean = false
	}

	if win.win.Maxmized {
		win.bounds = FourV2i{
			TopLeft:     V2i{X: 0, Y: 0},
			TopRight:    V2i{X: viewerWidth, Y: 0},
			BottomLeft:  V2i{X: 0, Y: viewerHeight},
			BottomRight: V2i{X: viewerWidth, Y: viewerHeight},
		}
	} else {
		win.bounds = FourV2i{
			TopLeft:     V2i{X: win.position.X, Y: win.position.Y},
			TopRight:    V2i{X: win.size.X + win.position.X, Y: win.position.Y},
			BottomLeft:  V2i{X: win.position.X, Y: win.size.Y + win.position.Y},
			BottomRight: V2i{X: win.size.X + win.position.X, Y: win.size.Y + win.position.Y},
		}
	}
	if win.win.HasTitleBar {
		win.titleBounds = FourV2i{
			TopLeft:     win.bounds.TopLeft,
			TopRight:    win.bounds.TopRight,
			BottomLeft:  V2i{X: win.bounds.TopLeft.X, Y: win.bounds.TopLeft.Y + win.win.TitleSize},
			BottomRight: V2i{X: win.bounds.BottomRight.X, Y: win.bounds.TopLeft.Y + win.win.TitleSize},
		}
		win.closeBounds = FourV2i{
			TopLeft:     V2i{X: win.bounds.TopRight.X - win.win.TitleSize, Y: win.bounds.TopRight.Y},
			TopRight:    win.bounds.TopRight,
			BottomLeft:  V2i{X: win.bounds.TopRight.X - win.win.TitleSize, Y: win.bounds.TopLeft.Y + win.win.TitleSize},
			BottomRight: V2i{X: win.bounds.TopRight.X, Y: win.bounds.TopLeft.Y + win.win.TitleSize},
		}
	}
}

func (rect FourV2i) contains(pos V2i) bool {
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
	var window *windowObject

	for _, win := range openWindows {
		if win.id == windowID {
			window = win
			break
		}
	}

	if window != nil {
		if window.open {
			window.open = false

			numOpen := len(openWindows) - 1
			if numOpen == 0 {
				openWindows = []*windowObject{}
				return nil
			}
			for w := numOpen; numOpen > 0; numOpen-- {
				if openWindows[w].id != windowID {
					openWindows = append(openWindows[:w], openWindows[w+1:]...)
					return nil
				}
			}
		}

		return nil
	}

	return errors.New("unable to find window")
}

func (win *windowObject) redraw() {
	if !win.clean {
		win.clean = true
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

	for _, win := range windowList {
		win.updateWin()
	}
}
