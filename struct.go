package seGUI

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var DefaultWinSettings = WindowData{
	HasTitleBar: true, Closable: true, Movable: true,
	TitleColor:       color.RGBA{R: 255, G: 255, B: 255, A: 255},
	TitleBGColor:     color.RGBA{R: 32, G: 32, B: 32, A: 255},
	TitleButtonColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
	BGColor:          color.RGBA{R: 16, G: 16, B: 16, A: 255},
}

type windowObject struct {
	id string

	win WindowData

	position,
	size V2i

	items []WindowItemData

	drawCache *ebiten.Image

	open,
	dirty bool
}

type WindowData struct {
	Title string

	StartPosition,
	StartSize V2i

	Closable, Focused, AutoCentered,
	Borderless, Movable, CachePersist,
	Resizable, KeepPosition, HasTitleBar bool

	TitleColor,
	TitleBGColor,
	TitleButtonColor,
	BGColor,
	BorderColor color.Color
}

type WindowItemData struct {
	Text     string
	Size     V2i
	Position V2i

	Color, HoverColor, ActionColor color.Color

	Action func()
}

type V2i struct {
	X, Y int
}
