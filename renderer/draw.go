package renderer

import (
	"github.com/faiface/pixel"
)

func (g *Renderer) drawText() {
	for row := 0; row < g.buffer.Size().Y; row++ {
		line := g.buffer.Row(row + g.buffer.ScrollY())
		if len(line.Chars) == 0 {
			continue
		}

		h := g.pixelHeight - (float64(row) * g.fontManager.CharSize().Y)
		g.text.DrawLine(g.window, line, pixel.V(0, h))
	}
}
