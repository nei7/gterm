package rect

import (
	"image/color"

	"github.com/faiface/pixel"
)

type Rect struct {
	frame  pixel.Rect
	d      pixel.Drawer
	tri    *pixel.TrianglesData
	matrix pixel.Matrix
	mask   pixel.RGBA
	pic    pixel.Picture
}

func NewRect(size pixel.Rect, color color.Color) *Rect {
	pic := pixel.MakePictureData(size)

	r, g, b, a := color.RGBA()
	for i := range pic.Pix {
		pic.Pix[i].R = uint8(r)
		pic.Pix[i].G = uint8(g)
		pic.Pix[i].B = uint8(b)
		pic.Pix[i].A = uint8(a)
	}

	tri := pixel.MakeTrianglesData(6)
	rect := &Rect{
		tri:  tri,
		d:    pixel.Drawer{Triangles: tri},
		mask: pixel.Alpha(1),
		pic:  pic,
	}

	rect.matrix = pixel.IM

	rect.d.Picture = pic

	rect.frame = pic.Rect
	rect.calcData()

	return rect
}

func (rect *Rect) Picture() pixel.Picture {
	return rect.pic
}

func (r *Rect) calcData() {
	var (
		center     = r.frame.Center()
		horizontal = pixel.V(r.frame.W()/2, 0)
		vertical   = pixel.V(0, r.frame.H()/2)
	)

	(*r.tri)[0].Position = pixel.Vec{}.Sub(horizontal).Sub(vertical)
	(*r.tri)[1].Position = pixel.Vec{}.Add(horizontal).Sub(vertical)
	(*r.tri)[2].Position = pixel.Vec{}.Add(horizontal).Add(vertical)
	(*r.tri)[3].Position = pixel.Vec{}.Sub(horizontal).Sub(vertical)
	(*r.tri)[4].Position = pixel.Vec{}.Add(horizontal).Add(vertical)
	(*r.tri)[5].Position = pixel.Vec{}.Sub(horizontal).Add(vertical)

	for i := range *r.tri {
		(*r.tri)[i].Color = r.mask
		(*r.tri)[i].Picture = center.Add((*r.tri)[i].Position)
		(*r.tri)[i].Intensity = 1
		(*r.tri)[i].Position = r.matrix.Project((*r.tri)[i].Position)
	}

	r.d.Dirty()
}

func (r *Rect) DrawColorMask(t pixel.Target, matrix pixel.Matrix, mask color.Color) {
	dirty := false
	if matrix != r.matrix {
		r.matrix = matrix
		dirty = true
	}
	if mask == nil {
		mask = pixel.Alpha(1)
	}
	rgba := pixel.ToRGBA(mask)
	if rgba != r.mask {
		r.mask = rgba
		dirty = true
	}

	if dirty {
		r.calcData()
	}

	r.d.Draw(t)
}

func DrawRect(t pixel.Target, mat pixel.Matrix, size pixel.Rect, color color.Color) *Rect {
	rect := NewRect(size, color)
	rect.DrawColorMask(t, mat, color)
	return rect
}
