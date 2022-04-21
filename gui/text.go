package gui

import (
	"image/color"
	"io/ioutil"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
	"github.com/goki/freetype/truetype"
	"github.com/nei7/gterm/term"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
)

// https://github.com/faiface/pixel/blob/master/text/text.go
type Text struct {
	Orig pixel.Vec
	Dot  pixel.Vec

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
}

type selection struct {
	bgColor color.Color
	mat     pixel.Matrix
}

func NewText(orig pixel.Vec, atlas *text.Atlas) *Text {
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

	size := txt.atlas.Glyph('Q')

	txt.batch, txt.sprite = createBatch(pixel.V(size.Frame.W(), size.Frame.H()+txt.LineHeight/4))

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

func loadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
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

func (txt *Text) Clear() {
	txt.prevR = -1
	txt.bounds = pixel.Rect{}
	txt.tris.SetLen(0)
	txt.dirty = true
	txt.Dot = txt.Orig
}

func (txt *Text) Draw(t pixel.Target, matrix pixel.Matrix, lines []term.Line) {
	txt.batch.Clear()

	txt.drawBuff(lines, t)

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

func (txt *Text) drawBuff(lines []term.Line, t pixel.Target) {

	txt.Clear()

	for _, line := range lines {
		for _, ch := range line.Chars {

			if ch.FgColor == nil {
				ch.FgColor = colornames.White
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

			if ch.BgColor != nil {
				w := frame.W()/txt.sprite.Frame().W() + txt.LineHeight/10

				moved := pixel.V(rect.Max.X-(rect.W()/2), txt.Dot.Y+txt.LineHeight/4)
				mat := pixel.IM.ScaledXY(pixel.ZV, pixel.V(w, 1)).Moved(moved)

				txt.sprite.DrawColorMask(txt.batch, mat, ch.BgColor)
			}

			if txt.bounds.W()*txt.bounds.H() == 0 {
				txt.bounds = bounds
			} else {
				txt.bounds = txt.bounds.Union(bounds)
			}

		}

		txt.Dot.X = txt.Orig.X
		txt.Dot.Y -= txt.LineHeight

	}

	txt.batch.Draw(t)
}
