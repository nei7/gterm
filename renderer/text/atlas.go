package text

import (
	"unicode"

	"github.com/faiface/pixel"
)

func (txt *Text) DrawLetter(r rune, pos pixel.Vec) (rect, frame pixel.Rect) {
	a := txt.Atlas()

	if !a.Contains(r) {
		r = unicode.ReplacementChar
	}
	if !a.Contains(unicode.ReplacementChar) {
		return pixel.Rect{}, pixel.Rect{}
	}

	txt.manager.RegularFontFace()
	glyph := a.Glyph(r)

	rect = glyph.Frame.Moved(pos.Sub(glyph.Dot))

	return rect, glyph.Frame
}
