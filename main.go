package main

import (
	"log"

	"github.com/faiface/pixel/pixelgl"
	"github.com/nei7/gterm/config"
	"github.com/nei7/gterm/renderer"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	t := renderer.NewRenderer(config)

	pixelgl.Run(t.Run)
}
