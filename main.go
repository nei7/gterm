package main

import (
	"log"

	"github.com/faiface/pixel/pixelgl"
	"github.com/nei7/gterm/config"
	"github.com/nei7/gterm/gui"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	t := gui.New(config)

	pixelgl.Run(t.Run)
}
