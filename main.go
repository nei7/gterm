package main

import (
	"log"

	"github.com/faiface/pixel/pixelgl"
)

func main() {

	t := NewTerminal()

	err := t.startPty()
	if err != nil {
		log.Fatal(err)
	}

	pixelgl.Run(t.runUI)

}
