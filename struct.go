package seGUI

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

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
