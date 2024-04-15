package EbitUI

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Run this in ebiten draw(), pass "screen"
func drawWindows(screen *ebiten.Image) {

	windowsLock.Lock()
	defer windowsLock.Unlock()

	for _, win := range openWindows {

		win.redraw()

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
