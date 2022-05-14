package text

import (
	"unicode"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"github.com/nei7/gterm/font"
	"github.com/nei7/gterm/renderer/rect"
	"github.com/nei7/gterm/term"
	"golang.org/x/image/colornames"
)

var charset = append(text.ASCII, text.RangeTable(unicode.Latin)...)

type Text struct {
	bold    *text.Atlas
	regular *text.Atlas
	prevR   rune
	trans   pixel.TrianglesData

	dw       pixel.Drawer
	dot      pixel.Vec
	charSize pixel.Vec
}

func NewText(manager *font.Manager) *Text {
	txt := &Text{
		regular:  text.NewAtlas(manager.RegularFontFace(), charset),
		bold:     text.NewAtlas(manager.BoldFontFace(), charset),
		charSize: manager.CharSize(),
	}

	pic := pixel.MakePictureData(pixel.R(0, 0, txt.charSize.X, txt.charSize.Y))

	r, g, b, a := colornames.Red.RGBA()
	for i := range pic.Pix {
		pic.Pix[i].R = uint8(r)
		pic.Pix[i].G = uint8(g)
		pic.Pix[i].B = uint8(b)
		pic.Pix[i].A = uint8(a)
	}

	txt.dw.Picture = txt.regular.Picture()
	txt.dw.Triangles = &txt.trans

	return txt
}

func (txt *Text) DrawLine(t pixel.Target, line term.Line, pos pixel.Vec, b *rect.RectBatch) {
	txt.trans.SetLen(0)
	var td pixel.TrianglesData
	td.SetLen(6)

	txt.dot = pos

	for _, ch := range line.Chars {

		if ch.BgColor == nil {
			ch.BgColor = colornames.Black
		}

		rectPos := pixel.IM.Moved(pixel.V(txt.dot.X+(txt.charSize.X/2), txt.dot.Y+txt.regular.Descent()))
		b.DrawColorMask(rectPos, ch.BgColor)

		var r, frame pixel.Rect
		r, frame, _, txt.dot = txt.regular.DrawRune(txt.prevR, ch.R, txt.dot)
		txt.prevR = ch.R

		rv := r.Vertices()
		fv := frame.Vertices()

		for i, j := range []int{0, 1, 2, 0, 2, 3} {
			td[i].Position = rv[j]
			td[i].Picture = fv[j]
		}

		if ch.FgColor == nil {
			ch.FgColor = colornames.White
		}

		for i := range td {
			td[i].Color = pixel.ToRGBA(ch.FgColor)
			td[i].Intensity = 1
		}

		txt.trans = append(txt.trans, td...)
		txt.dw.Dirty()

	}

	b.Draw(t)
	txt.dw.Draw(t)

	b.Clear()
}
