package utils

import (
	"io/ioutil"
	"os"

	"github.com/goki/freetype/truetype"
	"golang.org/x/image/font"
)


func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return home
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
