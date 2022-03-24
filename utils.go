package main

import (
	"io/ioutil"
	"os"

	"github.com/goki/freetype/truetype"
	"golang.org/x/image/font"
)

func loadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}

func calcOffset(len, offset, screenHeight, fontSize int) (startOffset int, endOffset int) {
	if offset > len {
		startOffset = len
	} else {
		startOffset = offset
	}
	if offset > 0 {

		endOffset = offset + (int(screenHeight) / (fontSize + 3)) - 3
		if endOffset > len {
			endOffset = len
			return
		}
		return
	}
	endOffset = len
	return
}
