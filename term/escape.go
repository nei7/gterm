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
	'P': escapeDeleteChars,
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
	rows, _ := strconv.Atoi(msg)
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
	t.moveCursor(t.buffer.savedCursorPos.Y, t.buffer.savedCursorPos.X)
}

func escapeSaveCursor(t *Terminal, _ string) {
	t.buffer.savedCursorPos.Y = t.buffer.cursorPos.Y
	t.buffer.savedCursorPos.X = t.buffer.cursorPos.X
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

func escapeEraseInScreen(t *Terminal, msg string) {
	mode, _ := strconv.Atoi(msg)
	switch mode {
	case 0:
		line := t.buffer.Row(t.buffer.cursorPos.Y)
		line.Chars = line.Chars[t.buffer.cursorPos.X:]
		t.buffer.lines = t.buffer.lines[t.buffer.cursorPos.Y:]
	case 1:
		line := t.buffer.Row(t.buffer.cursorPos.Y)
		line.Chars = line.Chars[:t.buffer.cursorPos.X]
		t.buffer.lines = t.buffer.lines[:t.buffer.cursorPos.Y]
	case 2:
		t.Clear()
	}
}

func escapeInsertLines(t *Terminal, msg string) {
	rows, _ := strconv.Atoi(msg)
	if rows == 0 {
		rows = 1
	}
	i := t.scrollOffset
	for ; i > t.buffer.cursorPos.Y-rows; i-- {
		t.buffer.SetRow(i, t.buffer.Row(i-rows).Chars)
	}
	for ; i >= t.buffer.cursorPos.Y; i-- {
		t.buffer.SetRow(i, []Char{})
	}
}

func escapeEraseInLine(t *Terminal, msg string) {
	mode, _ := strconv.Atoi(msg)
	switch mode {
	case 0:
		row := t.buffer.Row(t.buffer.cursorPos.Y)
		if t.buffer.cursorPos.X >= len(row.Chars) {
			return
		}
		t.buffer.SetRow(t.buffer.cursorPos.Y, row.Chars[:t.buffer.cursorPos.X])
	case 1:
		row := t.buffer.Row(t.buffer.cursorPos.Y)
		if t.buffer.cursorPos.X >= len(row.Chars) {
			return
		}
		chars := make([]Char, t.buffer.cursorPos.X)
		t.buffer.SetRow(t.buffer.cursorPos.Y, append(chars, row.Chars[t.buffer.cursorPos.X:]...))
	case 2:
		row := t.buffer.Row(t.buffer.cursorPos.Y)
		if t.buffer.cursorPos.X >= len(row.Chars) {
			return
		}
		chars := make([]Char, len(row.Chars))

		t.buffer.SetRow(t.buffer.cursorPos.Y, chars)
	}
}

func escapeDeleteChars(t *Terminal, msg string) {
	i, _ := strconv.Atoi(msg)
	right := t.buffer.cursorPos.X + i

	row := t.buffer.Row(t.buffer.cursorPos.Y)

	cells := row.Chars[:t.buffer.cursorPos.X]
	cells = append(cells, make([]Char, i)...)
	if right < len(t.buffer.lines) {
		cells = append(cells, row.Chars[right:]...)
	}

	t.buffer.SetRow(t.buffer.cursorPos.Y, cells)
}
