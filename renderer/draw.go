package renderer

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/nei7/gterm/renderer/rect"
)

func (g *Renderer) drawText() {
	size := g.fontManager.CharSize()
	rect := rect.NewRectBatch(pixel.R(0, 0, size.X*1.2, size.Y), color.White)

	for row := 0; row < g.buffer.Size().Y; row++ {
		line := g.buffer.Row(row + g.buffer.ScrollY())
		if len(line.Chars) == 0 {
			continue
		}

		h := g.pixelHeight - (float64(row) * size.Y)

		g.text.DrawLine(g.window, line, pixel.V(0, h), rect)
	}

}
