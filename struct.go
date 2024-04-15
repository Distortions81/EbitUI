package seGUI

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type WindowID string

var DefaultWinSettings = WindowData{
	Title:       "Window",
	TitleSize:   24,
	HasTitleBar: true, Closable: true, Movable: true,
	Resizable:        true,
	TitleColor:       color.RGBA{R: 255, G: 255, B: 255, A: 255},
	TitleBGColor:     color.RGBA{R: 32, G: 32, B: 32, A: 255},
	TitleButtonColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
	BGColor:          color.RGBA{R: 16, G: 16, B: 16, A: 255},
}

type windowObject struct {
	id    string
	layer int

	win WindowData

	position,
	oldPos,
	size,
	oldSize V2i

	bounds,
	titleBounds,
	closeBounds FourV2i

	items        []WindowItemData
	selectedItem *WindowItemData

	drawCache *ebiten.Image

	open, focused, mouseOver,
	clean bool
}

type WindowData struct {
	Title     string
	TitleSize int

	Position,
	Size V2i

	Closable, AutoCentered,
	Borderless, Movable, CachePersist,
	Resizable, KeepPosition, HasTitleBar, Maxmized bool

	TitleColor,
	TitleBGColor,
	TitleButtonColor,
	BGColor,
	BorderColor color.Color
}

type WindowItemData struct {
	IType    ITYPE
	FlowData FlowDataType

	Disabled, Hidden bool

	Tooltip string

	Label    string
	Size     V2i
	Position V2i

	Color, HoverColor, ActionColor color.Color

	Action func()
}

type FlowDataType struct {
	Parent     *WindowItemData
	Dir        FLOW_DIR
	Scrollable FLOW_DIR
	Resizeable bool
	Children   []WindowItemData
}

type FourV2i struct {
	TopLeft, TopRight, BottomLeft, BottomRight V2i
}

type V2i struct {
	X, Y int
}

type V2f struct {
	X, Y float32
}

type V2f64 struct {
	X, Y float64
}
