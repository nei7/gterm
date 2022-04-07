package term

import (
	"fmt"
	"image/color"
	"log"

	"golang.org/x/image/colornames"
	"golang.org/x/image/font"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

type char struct {
	r       rune
	fgColor color.Color
	bgColor color.Color
}

type Terminal struct {
	window  *pixelgl.Window
	text    *text.Text
	font    font.Face
	content []char

	cursorPos int
	offsetY   int

	height float64
	width  float64

	title string
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

	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:     "gterm",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     true,
		Resizable: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	atlas := text.NewAtlas(face, text.ASCII)
	text := text.New(pixel.Vec{
		X: config.Window.Padding.X,
		Y: config.Window.Padding.Y,
	}, atlas)

	t.text = text
	t.window = win
	t.font = face

	windowSize := win.Bounds().Size()
	t.width = windowSize.X
	t.height = windowSize.Y

	return t
}

func (t *Terminal) Write() {

}

func (t *Terminal) Update() {

	for !t.window.Closed() {
		t.window.Clear(colornames.Black)

		t.window.Update()
	}
}

func (t *Terminal) Run() {

}
