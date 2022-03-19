package main

import (
	"log"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

type UI struct {
	mu     sync.Mutex
	window *pixelgl.Window
	text   *text.Text
}

func NewUI() *UI {
	w := createWindow()

	ui := &UI{window: w}
	ui.loadFont()

	return ui
}

func createWindow() *pixelgl.Window {
	cfg := pixelgl.WindowConfig{
		Title:     "gterm",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return win
}

func (ui *UI) loadFont() {
	face, err := loadTTF("/usr/share/fonts/TTF/Sauce Code Pro Nerd Font Complete.ttf", 20)
	if err != nil {
		face = basicfont.Face7x13
	}
	atlas := text.NewAtlas(face, text.ASCII)
	text := text.New(pixel.V(10, ui.window.Canvas().Bounds().H()-30), atlas)
	ui.text = text
}
