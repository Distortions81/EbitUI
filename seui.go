package seGUI

import (
	"errors"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	windowList  map[string]*windowObject
	openWindows []*windowObject

	windowsLock sync.Mutex
)

func init() {
	windowList = map[string]*windowObject{}
}

// Run this in ebiten draw(), pass "screen"
func DrawWindows(screen *ebiten.Image) {

	windowsLock.Lock()
	defer windowsLock.Unlock()

	for _, win := range openWindows {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(win.position.X), float64(win.position.Y))
		screen.DrawImage(win.drawCache, op)
	}
}

// Add a window. Returns true if added
func AddWindow(windowID string, window WindowData) error {
	windowsLock.Lock()
	defer windowsLock.Unlock()

	newWin := &windowObject{win: window, dirty: true}
	newWin.size = newWin.win.StartSize
	windowList[windowID] = newWin

	newWin.drawCache = ebiten.NewImage(newWin.size.X, newWin.size.Y)
	if newWin.drawCache == nil {
		return errors.New("unable to create window draw cache")
	}

	newWin.drawCache.Fill(color.White)
	return nil
}

// Delete a window. Returns true if deleted
func DeleteWindow(windowID string) error {
	windowsLock.Lock()
	defer windowsLock.Unlock()

	if windowList[windowID] != nil {
		delete(windowList, windowID)
		return nil
	}

	return errors.New("unable to find window")
}

// Update a window. Returns true if updated
func UpdateWindow(windowID string, window WindowData) error {
	windowsLock.Lock()
	defer windowsLock.Unlock()

	if windowList[windowID] != nil {
		windowList[windowID].win = window
		windowList[windowID].dirty = true
		return nil
	}

	return errors.New("unable to find window")
}

// Update window items. Returns true if updated
func UpdateWindowItems(windowID string, windowItems []WindowItemData) error {
	windowsLock.Lock()
	defer windowsLock.Unlock()

	if windowList[windowID] != nil {
		windowList[windowID].items = windowItems
		windowList[windowID].dirty = true
		return nil
	}

	return errors.New("unable to find window")
}

func OpenWindow(windowID string) error {
	windowsLock.Lock()
	defer windowsLock.Unlock()

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
