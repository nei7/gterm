package term

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/image/colornames"
	"golang.org/x/image/font"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	atlas "github.com/faiface/pixel/text"
	"github.com/nei7/gterm/pkg/term/text"
)

type Terminal struct {
	window  *pixelgl.Window
	text    *text.Text
	font    font.Face
	content []text.Char

	pty *os.File

	cursorPos int

	height float64
	width  float64

	title string

	config *Config
}

func New() *Terminal {
	t := &Terminal{}

	config, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	fontPath := fmt.Sprintf("/usr/share/fonts/TTF/%s.ttf", config.Font.Family)

	face, err := loadTTF(fontPath, config.Font.Size)
	if err != nil {
		log.Fatal(err)
	}

	atlas := atlas.NewAtlas(face, atlas.ASCII)
	text := text.New(pixel.Vec{
		X: config.Window.Padding.X,
		Y: t.height - config.Window.Padding.Y,
	}, atlas)

	t.config = config
	t.text = text
	t.font = face

	return t
}

func (t *Terminal) setupWindow(w *pixelgl.Window) {
	windowSize := w.Bounds().Size()
	t.width = windowSize.X
	t.height = windowSize.Y
	t.window = w
}

func (t *Terminal) input() {
	typed := t.window.Typed()

	if typed != "" {

		runes := []rune(typed)
		for _, r := range runes {
			t.text.Chars = append(t.content, text.Char{
				Id:      t.cursorPos,
				FgColor: colornames.Azure,
				R:       r,
			})
			t.text.DrawBuf()

		}

		t.cursorPos++
	}

}

func (t *Terminal) draw() {

	cols := int(t.width / t.text.LineHeight)
	rows := int(t.height / t.text.TabWidth)
	i := 0
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			if i >= len(t.content) {
				return
			}
			i++
		}
	}

}

func (t *Terminal) Run() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:     "gterm",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     true,
		Resizable: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	t.setupWindow(win)

	for !win.Closed() {
		win.Clear(colornames.Black)

		t.input()
		t.draw()

		t.text.Draw(t.window, pixel.IM)
		win.Update()
	}
}
