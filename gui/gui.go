package gui

import (
	"fmt"
	"log"
	"unicode"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/nei7/gterm/config"
	"github.com/nei7/gterm/term"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

type GUI struct {
	window *pixelgl.Window

	height float64
	width  float64

	font font.Face

	text   *Text
	config *config.Config

	terminal *term.Terminal
}

func NewGUI(config *config.Config) *GUI {
	g := &GUI{}

	fontPath := fmt.Sprintf("/usr/share/fonts/TTF/%s.ttf", config.Font.Family)

	face, err := loadTTF(fontPath, config.Font.Size)
	if err != nil {
		log.Fatal(err)
	}

	g.font = face
	g.config = config

	g.terminal = term.New()

	return g
}

func (g *GUI) setupWindow(w *pixelgl.Window) {
	windowSize := w.Bounds().Size()
	g.width = windowSize.X - g.config.Window.Padding.X
	g.height = windowSize.Y - g.config.Window.Padding.Y

	g.window = w

	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII, text.RangeTable(unicode.Latin))

	g.text = NewText(pixel.Vec{
		X: g.config.Window.Padding.X,
		Y: g.height - atlas.Ascent() - g.config.Window.Padding.Y,
	}, atlas)

	g.Resize()
}

func (g *GUI) Resize() {
	windowSize := g.window.Bounds().Size()

	cols := int(g.width / g.text.atlas.Glyph(' ').Advance)
	rows := int(g.height / (g.text.LineHeight))

	g.terminal.SetSize(uint16(rows), uint16(cols))

	g.terminal.SetPtySize(windowSize)

}
