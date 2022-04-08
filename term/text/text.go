package text

import (
	"image/color"
	"math"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
)

type Char struct {
	Id      int
	R       rune
	FgColor color.Color
	BgColor color.Color
}

// https://github.com/faiface/pixel/blob/master/text/text.go
type Text struct {
	sync.RWMutex

	chars [][]Char
	Orig  pixel.Vec

	cols         int
	rows         int
	scrollOffset int

	Dot pixel.Vec

	Color color.Color

	LineHeight float64

	TabWidth float64

	atlas *text.Atlas

	prevR  rune
	bounds pixel.Rect
	glyph  pixel.TrianglesData
	tris   pixel.TrianglesData

	mat    pixel.Matrix
	col    pixel.RGBA
	trans  pixel.TrianglesData
	transD pixel.Drawer
	dirty  bool
}

func New(orig pixel.Vec, atlas *text.Atlas) *Text {
	txt := &Text{
		Orig:       orig,
		Dot:        orig,
		Color:      pixel.Alpha(1),
		LineHeight: atlas.LineHeight(),
		TabWidth:   atlas.Glyph(' ').Advance * 4,
		atlas:      atlas,
		mat:        pixel.IM,
		col:        pixel.Alpha(1),
	}

	txt.glyph.SetLen(6)
	for i := range txt.glyph {
		txt.glyph[i].Color = pixel.Alpha(1)
		txt.glyph[i].Intensity = 1
	}

	txt.transD.Picture = txt.atlas.Picture()
	txt.transD.Triangles = &txt.trans

	txt.Clear()

	return txt
}

func (txt *Text) SetSize(rows, cols int) {
	txt.rows = rows
	txt.cols = cols
}

func (txt *Text) Write(buf []byte) {
	txt.Lock()
	defer txt.Unlock()

	txt.handleOutput(buf)
	txt.ScrollDown()
	txt.drawBuff()
}

func (txt *Text) Atlas() *text.Atlas {
	return txt.atlas
}

func (txt *Text) Bounds() pixel.Rect {
	return txt.bounds
}

func (txt *Text) BoundsOf(s string) pixel.Rect {
	dot := txt.Dot
	prevR := txt.prevR
	bounds := pixel.Rect{}

	for _, r := range s {
		var control bool
		dot, control = txt.controlRune(r, dot)
		if control {
			continue
		}

		var b pixel.Rect
		_, _, b, dot = txt.Atlas().DrawRune(prevR, r, dot)

		if bounds.W()*bounds.H() == 0 {
			bounds = b
		} else {
			bounds = bounds.Union(b)
		}

		prevR = r
	}

	return bounds
}

func (txt *Text) Scroll(speed int) {
	if txt.scrollOffset-speed >= 0 {

		if txt.scrollOffset-speed+txt.rows > len(txt.chars) {
			return
		}

		txt.scrollOffset -= speed
	}

	txt.drawBuff()
}

func (txt *Text) ScrollDown() {
	txt.scrollOffset = len(txt.chars) - txt.rows
	if txt.scrollOffset < 0 {
		txt.scrollOffset = 0
	}

}

func (txt *Text) Clear() {
	txt.prevR = -1
	txt.bounds = pixel.Rect{}
	txt.tris.SetLen(0)
	txt.dirty = true
	txt.Dot = txt.Orig
}

func (txt *Text) Draw(t pixel.Target, matrix pixel.Matrix) {
	txt.RLock()
	defer txt.RUnlock()

	if matrix != txt.mat {
		txt.mat = matrix
		txt.dirty = true
	}

	if txt.dirty {
		txt.trans.SetLen(txt.tris.Len())

		txt.trans.Update(&txt.tris)

		for i := range txt.trans {
			txt.trans[i].Position = txt.mat.Project(txt.trans[i].Position)
			txt.trans[i].Color = txt.trans[i].Color.Mul(txt.col)
		}

		txt.transD.Dirty()
		txt.dirty = false
	}

	txt.transD.Draw(t)
}

func (txt *Text) controlRune(r rune, dot pixel.Vec) (newDot pixel.Vec, control bool) {
	switch r {
	case '\r':
		dot.X = txt.Orig.X
	case '\t':
		rem := math.Mod(dot.X-txt.Orig.X, txt.TabWidth)
		rem = math.Mod(rem, rem+txt.TabWidth)
		if rem == 0 {
			rem = txt.TabWidth
		}
		dot.X += rem

	default:
		return dot, false
	}
	return dot, true
}

func (txt *Text) drawBuff() {

	txt.Clear()

	endOff := txt.scrollOffset + txt.rows
	if endOff > len(txt.chars) {
		endOff = len(txt.chars)
	}

	for _, lines := range txt.chars[txt.scrollOffset:endOff] {
		for _, ch := range lines {
			var control bool
			txt.Dot, control = txt.controlRune(ch.R, txt.Dot)
			if control {
				continue
			}

			for i := range txt.glyph {
				txt.glyph[i].Color = pixel.ToRGBA(ch.FgColor)
			}

			var rect, frame, bounds pixel.Rect
			rect, frame, bounds, txt.Dot = txt.Atlas().DrawRune(txt.prevR, ch.R, txt.Dot)

			txt.prevR = ch.R

			rv := [...]pixel.Vec{
				{X: rect.Min.X, Y: rect.Min.Y},
				{X: rect.Max.X, Y: rect.Min.Y},
				{X: rect.Max.X, Y: rect.Max.Y},
				{X: rect.Min.X, Y: rect.Max.Y},
			}

			fv := [...]pixel.Vec{
				{X: frame.Min.X, Y: frame.Min.Y},
				{X: frame.Max.X, Y: frame.Min.Y},
				{X: frame.Max.X, Y: frame.Max.Y},
				{X: frame.Min.X, Y: frame.Max.Y},
			}

			for i, j := range [...]int{0, 1, 2, 0, 2, 3} {
				txt.glyph[i].Position = rv[j]
				txt.glyph[i].Picture = fv[j]
			}

			txt.tris = append(txt.tris, txt.glyph...)
			txt.dirty = true

			if txt.bounds.W()*txt.bounds.H() == 0 {
				txt.bounds = bounds
			} else {
				txt.bounds = txt.bounds.Union(bounds)
			}
		}

		txt.Dot.X = txt.Orig.X
		txt.Dot.Y -= txt.LineHeight

	}

}
