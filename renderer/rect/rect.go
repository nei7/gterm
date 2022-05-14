package rect

import (
	"image/color"

	"github.com/faiface/pixel"
)

type RectBatch struct {
	frame pixel.Rect
	d     pixel.Drawer

	tris pixel.TrianglesData
	tmp  *pixel.TrianglesData

	matrix pixel.Matrix
	mask   pixel.RGBA
	pic    pixel.Picture
}

func NewRectBatch(size pixel.Rect, color color.Color) *RectBatch {
	pic := pixel.MakePictureData(size)

	r, g, b, a := color.RGBA()
	for i := range pic.Pix {
		pic.Pix[i].R = uint8(r)
		pic.Pix[i].G = uint8(g)
		pic.Pix[i].B = uint8(b)
		pic.Pix[i].A = uint8(a)
	}

	rect := &RectBatch{
		mask: pixel.Alpha(1),
		pic:  pic,
	}

	rect.d = pixel.Drawer{Triangles: &rect.tris, Picture: pic}

	rect.matrix = pixel.IM

	rect.frame = pic.Rect

	return rect
}

func (rect *RectBatch) Picture() pixel.Picture {
	return rect.pic
}

func (r *RectBatch) calcData() {
	var (
		center     = r.frame.Center()
		horizontal = pixel.V(r.frame.W()/2, 0)
		vertical   = pixel.V(0, r.frame.H()/2)
	)

	tri := *pixel.MakeTrianglesData(6)

	tri[0].Position = pixel.Vec{}.Sub(horizontal).Sub(vertical)
	tri[1].Position = pixel.Vec{}.Add(horizontal).Sub(vertical)
	tri[2].Position = pixel.Vec{}.Add(horizontal).Add(vertical)
	tri[3].Position = pixel.Vec{}.Sub(horizontal).Sub(vertical)
	tri[4].Position = pixel.Vec{}.Add(horizontal).Add(vertical)
	tri[5].Position = pixel.Vec{}.Sub(horizontal).Add(vertical)

	for i := range tri {
		tri[i].Color = r.mask
		tri[i].Picture = center.Add(tri[i].Position)
		tri[i].Intensity = 1
		tri[i].Position = r.matrix.Project(tri[i].Position)
	}

	r.tris = append(r.tris, tri...)

	r.d.Dirty()
}

func (r *RectBatch) DrawColorMask(matrix pixel.Matrix, mask color.Color) {
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
}

func (r *RectBatch) Draw(t pixel.Target) {
	r.d.Draw(t)
}

func (r *RectBatch) Clear() {
	r.d.Triangles.SetLen(0)
	r.d.Dirty()
}
