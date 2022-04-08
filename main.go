package main

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/nei7/gterm/pkg/term"
)

func main() {
	t := term.New()
	pixelgl.Run(t.Run)
}
