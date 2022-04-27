package term

import (
	"image/color"
	"log"
	"strconv"
	"strings"

	"golang.org/x/image/colornames"
)

var (
	basicColors = []color.RGBA{
		colornames.Black,
		{170, 0, 0, 255},
		{0, 170, 0, 255},
		{170, 170, 0, 255},
		{0, 0, 170, 255},
		{170, 0, 170, 255},
		{0, 255, 255, 255},
		{170, 170, 170, 255},
	}
	brightColors = []color.RGBA{
		{85, 85, 85, 255},
		{255, 85, 85, 255},
		{85, 255, 85, 255},
		{255, 255, 85, 255},
		{85, 85, 255, 255},
		{255, 85, 255, 255},
		{85, 255, 255, 255},
		{255, 255, 255, 255},
	}
)
var escapes = map[rune]func(*Terminal, string){
	'A': escapeMoveCursorUp,
	'B': escapeMoveCursorDown,
	'C': escapeMoveCursorRight,
	'D': escapeMoveCursorLeft,
	'd': escapeMoveCursorRow,
	'H': escapeMoveCursor,
	'f': escapeMoveCursor,
	'G': escapeMoveCursorCol,
	'L': escapeInsertLines,
	'm': escapeColorMode,
	'J': escapeEraseInScreen,
	'K': escapeEraseInLine,
	'r': escapeSetScrollArea,
	's': escapeSaveCursor,
	'u': escapeRestoreCursor,
}

func (t *Terminal) handleEscape(code string) {
	code = trimLeftZeros(code)
	if code == "" {
		return
	}

	runes := []rune(code)

	if esc, ok := escapes[runes[len(code)-1]]; ok {
		esc(t, code[:len(code)-1])
	} else if t.debug {
		log.Println("Unrecognised Escape:", code)
	}
}

func trimLeftZeros(s string) string {
	if s == "" {
		return s
	}

	i := 0
	for _, r := range s {
		if r > '0' {
			break
		}
		i++
	}

	return s[i:]
}

func (t *Terminal) moveCursor(row, col int) {

	t.buffer.cursorPos.X = col
	t.buffer.cursorPos.Y = row

}

func escapeMoveCursorUp(t *Terminal, msg string) {
	rows, err := strconv.Atoi(msg)
	if err != nil {
		log.Fatal(err)
	}

	if rows == 0 {
		rows = 1
	}

	t.moveCursor(t.buffer.cursorPos.Y-rows, t.buffer.cursorPos.X)
}

func escapeMoveCursorDown(t *Terminal, msg string) {
	rows, _ := strconv.Atoi(msg)
	if rows == 0 {
		rows = 1
	}

	t.moveCursor(t.buffer.cursorPos.Y+rows, t.buffer.cursorPos.X)
}

func escapeMoveCursorRight(t *Terminal, msg string) {
	cols, _ := strconv.Atoi(msg)
	if cols == 0 {
		cols = 1
	}
	t.moveCursor(t.buffer.cursorPos.Y, t.buffer.cursorPos.X+cols)
}

func escapeMoveCursorLeft(t *Terminal, msg string) {
	cols, _ := strconv.Atoi(msg)
	if cols == 0 {
		cols = 1
	}

	t.moveCursor(t.buffer.cursorPos.Y, t.buffer.cursorPos.X-cols)
}

func escapeMoveCursorRow(t *Terminal, msg string) {
	row, _ := strconv.Atoi(msg)
	t.moveCursor(row-1, t.buffer.cursorPos.X)
}

func escapeMoveCursorCol(t *Terminal, msg string) {
	col, _ := strconv.Atoi(msg)
	t.moveCursor(t.buffer.cursorPos.Y, col-1)
}

func escapeMoveCursor(t *Terminal, msg string) {
	if !strings.Contains(msg, ";") {
		t.moveCursor(0, 0)
		return
	}

	parts := strings.Split(msg, ";")
	row, _ := strconv.Atoi(parts[0])
	col := 1
	if len(parts) == 2 {
		col, _ = strconv.Atoi(parts[1])
	}

	t.moveCursor(row-1, col-1)
}

func escapeRestoreCursor(t *Terminal, _ string) {
	t.moveCursor(int(t.buffer.savedRows), int(t.buffer.savedCols))
}

func escapeSaveCursor(t *Terminal, _ string) {
	t.buffer.savedRows = uint16(t.buffer.cursorPos.Y)
	t.buffer.savedCols = uint16(t.buffer.cursorPos.X)
}

func escapeSetScrollArea(t *Terminal, msg string) {
	parts := strings.Split(msg, ";")
	start := 0
	end := int(t.buffer.rows) - 1
	if len(parts) == 2 {
		if parts[0] != "" {
			start, _ = strconv.Atoi(parts[0])
			start--
		}
		if parts[1] != "" {
			end, _ = strconv.Atoi(parts[1])
			end--
		}
	}
}

func escapeColorMode(t *Terminal, msg string) {
	if msg == "" || msg == "0" {
		t.currentBG = nil
		t.currentFG = nil
		return
	}

	modes := strings.Split(msg, ";")
	mode := modes[0]
	if (mode == "38" || mode == "48") && len(modes) >= 2 {
		if modes[1] == "5" && len(modes) >= 3 {
			t.handleColorModeMap(mode, modes[2])
			modes = modes[3:]
		} else if modes[1] == "2" && len(modes) >= 5 {
			t.handleColorModeRGB(mode, modes[2], modes[3], modes[4])
			modes = modes[5:]
		}
	}
	for _, mode := range modes {
		if mode == "" {
			continue
		}
		t.handleColorMode(mode)
	}
}

func (t *Terminal) handleColorModeRGB(mode, rs, gs, bs string) {
	r, _ := strconv.Atoi(rs)
	g, _ := strconv.Atoi(gs)
	b, _ := strconv.Atoi(bs)
	c := &color.RGBA{uint8(r), uint8(g), uint8(b), 255}

	if mode == "38" {
		t.currentFG = c
	} else if mode == "48" {
		t.currentBG = c
	}
}

func (t *Terminal) handleColorMode(modeStr string) {
	mode, err := strconv.Atoi(modeStr)
	if err != nil {
		log.Fatal(err)
	}
	switch mode {
	case 0:
		t.currentBG, t.currentFG = nil, nil
		t.bright = false
	case 1:
		t.bright = true
	case 4, 24: // italic
	case 7: // reverse
		bg := t.currentBG
		if t.currentFG == nil {
			t.currentBG = &colornames.White
		} else {
			t.currentBG = t.currentFG
		}
		if bg == nil {
			t.currentFG = &colornames.White
		} else {
			t.currentFG = bg
		}
	case 27: // reverse off
		bg := t.currentBG
		if t.currentFG == &colornames.White {
			t.currentBG = nil
		} else {
			t.currentBG = t.currentFG
		}
		if bg == &colornames.Black {
			t.currentFG = nil
		} else {
			t.currentFG = bg
		}
	case 30, 31, 32, 33, 34, 35, 36, 37:
		if t.bright {
			t.currentFG = brightColors[mode-30]
		} else {
			t.currentFG = basicColors[mode-30]
		}
	case 39:
		t.currentFG = nil
	case 40, 41, 42, 43, 44, 45, 46, 47:
		if t.bright {
			t.currentBG = brightColors[mode-40]
		} else {
			t.currentBG = basicColors[mode-40]
		}
	case 49:
		t.currentBG = nil
	case 90, 91, 92, 93, 94, 95, 96, 97:
		t.currentFG = brightColors[mode-90]
	case 100, 101, 102, 103, 104, 105, 106, 107:
		t.currentBG = brightColors[mode-100]
	default:
		if t.debug {
			log.Println("Unsupported graphics mode", mode)
		}
	}
}

func (t *Terminal) handleColorModeMap(mode, ids string) {
	var c color.Color
	id, err := strconv.Atoi(ids)
	if err != nil {
		if t.debug {
			log.Println("Invalid color map ID", ids)
		}
		return
	}
	if id <= 7 {
		c = basicColors[id]
	} else if id <= 15 {
		c = brightColors[id-8]
	} else if id <= 231 {
		inc := 256 / 5
		id -= 16
		b := id % 6
		id = (id - b) / 6
		g := id % 6
		r := (id - g) / 6
		c = &color.RGBA{uint8(r * inc), uint8(g * inc), uint8(b * inc), 255}
	} else if id <= 255 {
		id -= 232
		inc := 256 / 24
		y := id * inc
		c = &color.Gray{uint8(y)}
	} else if t.debug {
		log.Println("Invalid colour map ID", id)
	}

	if mode == "38" {

		t.currentFG = c
	} else if mode == "48" {

		t.currentBG = c
	}
}

func escapeEraseInScreen(t *Terminal, msg string) {
	mode, _ := strconv.Atoi(msg)
	switch mode {
	case 0:

	case 1:

	case 2:
		t.Clear()
	}
}

func escapeInsertLines(t *Terminal, msg string) {
	rows, _ := strconv.Atoi(msg)
	if rows == 0 {
		rows = 1
	}

}

func escapeEraseInLine(t *Terminal, msg string) {
	mode, _ := strconv.Atoi(msg)
	switch mode {
	case 0:
		line := t.buffer.getLine(t.buffer.cursorPos.Y)
		line.Chars = line.Chars[:t.buffer.cursorPos.X]
	case 1:
		line := t.buffer.getLine(t.buffer.cursorPos.Y)
		line.Chars = line.Chars[t.buffer.cursorPos.X:]
	case 2:
		t.buffer.clear()
	}
}
