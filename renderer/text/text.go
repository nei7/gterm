package text

import (
	"image/color"
	"unicode"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	manager "github.com/nei7/gterm/font"
	"github.com/nei7/gterm/term"
	"golang.org/x/image/colornames"
)

// https://github.com/faiface/pixel/blob/master/text/text.go
type Text struct {
	Orig pixel.Vec

	Color color.Color

	LineHeight float64
	TabWidth   float64

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

	batch  *pixel.Batch
	sprite *pixel.Sprite

	manager *manager.Manager
}

type selection struct {
	bgColor color.Color
	mat     pixel.Matrix
}

func NewText(orig pixel.Vec, fontManager *manager.Manager) *Text {

	atlas := text.NewAtlas(fontManager.RegularFontFace(), text.ASCII, text.RangeTable(unicode.Latin))

	txt := &Text{
		Orig:       orig,
		Color:      pixel.Alpha(1),
		LineHeight: atlas.LineHeight(),
		TabWidth:   atlas.Glyph(' ').Advance * 4,
		atlas:      atlas,

		mat:     pixel.IM,
		col:     pixel.Alpha(1),
		manager: fontManager,
	}

	txt.glyph.SetLen(6)
	for i := range txt.glyph {
		txt.glyph[i].Color = pixel.Alpha(1)
		txt.glyph[i].Intensity = 1
	}

	txt.transD.Picture = txt.atlas.Picture()
	txt.transD.Triangles = &txt.trans

	txt.Clear()

	size := fontManager.CharSize()
	txt.batch, txt.sprite = createBatch(pixel.V(float64(size.X), float64(size.Y)))

	return txt
}

func createBatch(size pixel.Vec) (*pixel.Batch, *pixel.Sprite) {
	spritesheet := createSpritesheet(size, colornames.White)
	batch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
	sprite := pixel.NewSprite(spritesheet, spritesheet.Rect)

	return batch, sprite
}

func createSpritesheet(size pixel.Vec, color color.Color) *pixel.PictureData {
	rect := pixel.R(0, 0, size.X, size.Y)

	r, g, b, a := color.RGBA()
	spritesheet := pixel.MakePictureData(rect)
	for i := range spritesheet.Pix {
		spritesheet.Pix[i].R = uint8(r)
		spritesheet.Pix[i].G = uint8(g)
		spritesheet.Pix[i].B = uint8(b)
		spritesheet.Pix[i].A = uint8(a)
	}

	return spritesheet
}

func (txt *Text) Atlas() *text.Atlas {
	return txt.atlas
}

func (txt *Text) Bounds() pixel.Rect {
	return txt.bounds
}

func (txt *Text) Clear() {
	txt.prevR = -1
	txt.bounds = pixel.Rect{}
	txt.tris.SetLen(0)
	txt.dirty = true
}

func (txt *Text) Draw(t pixel.Target, lines []term.Line) {
	txt.Clear()
	txt.batch.Clear()
	txt.drawBuff(lines, t)

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

func (txt *Text) drawBuff(lines []term.Line, t pixel.Target) {

	for y, line := range lines {
		for x, ch := range line.Chars {

			if ch.FgColor == nil {
				ch.FgColor = color.White
			}

			for i := range txt.glyph {
				txt.glyph[i].Color = pixel.ToRGBA(ch.FgColor)
			}

			size := txt.manager.CharSize()
			pos := pixel.V((float64(size.X) * float64(x)), txt.Orig.Y-float64(size.Y)*float64(y))

			rect, frame := txt.DrawLetter(ch.R, pos)

			pos.X += size.X / 2
			pos.Y += txt.atlas.Descent()

			if ch.BgColor != nil {
				txt.sprite.DrawColorMask(txt.batch, pixel.IM.Moved(pos), ch.BgColor)
			}

			txt.prevR = ch.R

			rv := []pixel.Vec{
				{X: rect.Min.X, Y: rect.Min.Y},
				{X: rect.Max.X, Y: rect.Min.Y},
				{X: rect.Max.X, Y: rect.Max.Y},
				{X: rect.Min.X, Y: rect.Max.Y},
			}

			fv := []pixel.Vec{
				{X: frame.Min.X, Y: frame.Min.Y},
				{X: frame.Max.X, Y: frame.Min.Y},
				{X: frame.Max.X, Y: frame.Max.Y},
				{X: frame.Min.X, Y: frame.Max.Y},
			}

			for i, j := range []int{0, 1, 2, 0, 2, 3} {
				txt.glyph[i].Position = rv[j]
				txt.glyph[i].Picture = fv[j]
			}

			txt.tris = append(txt.tris, txt.glyph...)
		}
	}

	txt.batch.Draw(t)
}
