package seGUI

import (
	"bytes"
	"errors"
	"image/color"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Init, with starting screen width and height
func Start(width, height int) {
	UpdateViewerSize(width, height)

	windowsLock.Lock()

	windowList = map[WindowID]*windowObject{}

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

	newWin.updateWin()
	windowList[WindowID(windowID)] = newWin

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

	if windowList[WindowID(windowID)] != nil {
		delete(windowList, WindowID(windowID))
		return nil
	}

	return errors.New("unable to find window")
}

func OpenWindow(windowID string) error {
	windowsLock.Lock()
	defer windowsLock.Unlock()
	windowID = strings.ToLower(windowID)

	window := windowList[WindowID(windowID)]

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

	return closeWindow(windowID)
}

func UpdateViewerSize(width, height int) (int, int) {
	if width < minSizeX {
		width = minSizeX
	}
	if height < minSizeY {
		height = minSizeY
	}

	clampWindows()
	return width, height
}

// Run this in ebiten draw(), pass "screen"
func DrawWindows(screen *ebiten.Image) {
	drawWindows(screen)
}
