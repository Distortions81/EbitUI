package seGUI

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"log"
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
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

func GetBounds(win *windowObject) FourV2i {
	halfx, halfy := win.size.X/2, win.size.Y/2
	rect := FourV2i{
		TopLeft:     V2i{X: win.position.X - halfx, Y: win.position.Y - halfy},
		TopRight:    V2i{X: win.position.X + halfx, Y: win.position.Y - halfy},
		BottomLeft:  V2i{X: win.position.X - halfx, Y: win.position.Y + halfy},
		BottomRight: V2i{X: win.position.X + halfx, Y: win.position.Y + halfy},
	}
	return rect
}

func UpdateWinPos(win *windowObject, pos V2i) {
	win.position = pos
	win.bounds = GetBounds(win)
}

func UpdateWinSize(win *windowObject, size V2i) {
	win.size = size
	win.bounds = GetBounds(win)
}

// Input update, returns if it ate: left click, right click, or ebiten.key
func InputUpdate() (bool, bool, int) {
	mx, my := ebiten.CursorPosition()

	//Detect clicks within open windows
	for _, item := range openWindows {
		if PosWithinRect(V2i{X: mx, Y: my}, item.rect) {
			return true, false, 0
		}
	}
	return false, false, 0
}

func PosWithinRect(pos V2i, rect FourV2i) bool {
	if pos.X >= rect.TopLeft.X &&
		pos.Y >= rect.TopLeft.Y &&
		pos.X <= rect.BottomRight.X &&
		pos.Y <= rect.BottomRight.Y {
		return true
	}
	return false
}

// Init, with starting screen width and height
func Start(width, height int) {
	UpdateViewerSize(width, height)

	windowsLock.Lock()

	windowList = map[string]*windowObject{}

	//Used for vectors
	whiteImage.Fill(color.White)

	//Load default font
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s

	nw := DefaultWinSettings
	nw.StartSize = V2i{X: width, Y: height}
	windowsLock.Unlock()

	err = AddWindow("hud", nw)
	if err != nil {
		log.Fatal(err)
	}
}

// Add a window. Returns true if added
func AddWindow(windowID string, window WindowData) error {
	windowsLock.Lock()
	defer windowsLock.Unlock()
	windowID = strings.ToLower(windowID)

	newWin := &windowObject{win: window, dirty: true}

	newWin.size = newWin.win.StartSize
	if window.HasTitleBar {
		newWin.size.Y += window.TitleSize
	}
	newWin.position = newWin.win.StartPosition
	windowList[windowID] = newWin

	newWin.drawCache = ebiten.NewImage(newWin.size.X, newWin.size.Y)
	if newWin.drawCache == nil {
		return errors.New("unable to create window draw cache")
	}

	newWin.drawCache.Fill(newWin.win.BGColor)
	return nil
}

// Delete a window. Returns true if deleted
func DeleteWindow(windowID string) error {
	windowsLock.Lock()
	defer windowsLock.Unlock()
	windowID = strings.ToLower(windowID)

	if windowList[windowID] != nil {
		delete(windowList, windowID)
		return nil
	}

	return errors.New("unable to find window")
}

func OpenWindow(windowID string) error {
	windowsLock.Lock()
	defer windowsLock.Unlock()
	windowID = strings.ToLower(windowID)

	window := windowList[windowID]

	if window != nil {
		if !window.open {
			window.open = true

			for _, win := range openWindows {
				if win.id == windowID {
					return nil
				}
			}
			openWindows = append(openWindows, window)
		}
		return nil
	}

	return errors.New("unable to find window")
}

func CloseWindow(windowID string) error {
	windowsLock.Lock()
	defer windowsLock.Unlock()
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

func UpdateViewerSize(width, height int) (int, int) {
	if width < minSizeX {
		width = minSizeX
	}
	if height < minSizeY {
		height = minSizeY
	}

	clampWindows(width, height)
	return width, height
}

// Run this in ebiten draw(), pass "screen"
func DrawWindows(screen *ebiten.Image) {
	drawWindows(screen)
}
