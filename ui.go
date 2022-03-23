package main

import (
	"log"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

const (
	fontSize = 15
)

var screenHeight float64

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

func (t *Terminal) loadFont() {
	face, err := loadTTF("/usr/share/fonts/TTF/Sauce Code Pro Nerd Font Complete.ttf", 15)
	if err != nil {
		face = basicfont.Face7x13
	}
	atlas := text.NewAtlas(face, text.ASCII)

	screenHeight = t.window.Canvas().Bounds().H()
	text := text.New(pixel.V(5, screenHeight-20), atlas)
	t.text = text
}
