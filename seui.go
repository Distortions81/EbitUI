package seGUI

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	screenWidth, screenHeight int
	windowList                map[string]*windowObject
	openWindows               []*windowObject

	windowsLock sync.Mutex

	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)

	mplusFaceSource *text.GoTextFaceSource
)

// Init, with starting screen width and height
func Start(width, height int) {
	UpdateScreenSize(width, height)

	windowsLock.Lock()
	defer windowsLock.Unlock()

	windowList = map[string]*windowObject{}

	//Used for vectors
	whiteImage.Fill(color.White)

	//Load default font
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
}

// Run this in ebiten draw(), pass "screen"
func DrawWindows(screen *ebiten.Image) {

	windowsLock.Lock()
	defer windowsLock.Unlock()

	for _, win := range openWindows {
		//Draw window
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(win.position.X), float64(win.position.Y))
		screen.DrawImage(win.drawCache, op)

		if !win.win.HasTitleBar {
			continue
		}

		//Title bg color
		vector.DrawFilledRect(screen, float32(win.position.X), float32(win.position.Y),
			float32(win.size.X), float32(win.win.TitleSize), win.win.TitleBGColor, false)

		//Title text
		loo := text.LayoutOptions{
			LineSpacing:    1,
			PrimaryAlign:   text.AlignStart,
			SecondaryAlign: text.AlignStart,
		}
		tdop := ebiten.DrawImageOptions{}
		tdop.GeoM.Translate(float64(win.position.X+2), float64(win.position.Y-2))

		top := &text.DrawOptions{DrawImageOptions: tdop, LayoutOptions: loo}
		text.Draw(screen, win.win.Title, &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   float64(win.win.TitleSize - 4),
		}, top)

		//Draw close X
		if win.win.Closable {

			tr := V2i{X: win.position.X + win.size.X - win.win.TitleSize/4, Y: win.position.Y + win.win.TitleSize/4}
			var path vector.Path
			path.MoveTo(float32(tr.X), float32(tr.Y))
			path.LineTo(float32(tr.X-win.win.TitleSize/2), float32(tr.Y+win.win.TitleSize/2))

			path.MoveTo(float32(tr.X-win.win.TitleSize/2), float32(tr.Y))
			path.LineTo(float32(tr.X), float32(tr.Y+win.win.TitleSize/2))

			path.Close()

			var vs []ebiten.Vertex
			var is []uint16
			vop := &vector.StrokeOptions{Width: 5, LineJoin: vector.LineJoinRound, LineCap: vector.LineCapRound}
			vs, is = path.AppendVerticesAndIndicesForStroke(nil, nil, vop)

			red, green, blue, alpha := win.win.TitleButtonColor.RGBA()
			for i := range vs {
				vs[i].ColorR = float32(red / 255)
				vs[i].ColorG = float32(green / 255)
				vs[i].ColorB = float32(blue / 255)
				vs[i].ColorA = float32(alpha / 255)
			}

			top := &ebiten.DrawTrianglesOptions{AntiAlias: true, FillRule: ebiten.FillAll}
			screen.DrawTriangles(vs, is, whiteSubImage, top)
		}
	}
}

// Add a window. Returns true if added
func AddWindow(windowID string, window WindowData) error {
	windowsLock.Lock()
	defer windowsLock.Unlock()

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

	if windowList[windowID] != nil {
		delete(windowList, windowID)
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

// Call his in Ebiten Layout
func clampWindows(width, height int) {
	windowsLock.Lock()
	defer windowsLock.Unlock()

	if width == screenWidth &&
		height == screenHeight {
		return
	}

	screenWidth = width
	screenHeight = height

	for _, win := range windowList {
		if win.size.X > screenWidth {
			win.size.X = screenWidth
		}
		if win.size.Y > screenHeight {
			win.size.Y = screenHeight
		}

		if win.position.X+win.size.X > screenWidth {
			win.position.X = (screenWidth - win.size.X)
		}
		if win.position.Y+win.size.Y > screenHeight {
			win.position.Y = (screenHeight - win.size.Y)
		}
	}
}

func UpdateScreenSize(width, height int) {
	clampWindows(width, height)
}
