package term

import (
	"log"
	"strconv"
	"strings"
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
	'r': escapeSetScrollArea,
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
	right, _ := strconv.Atoi(msg)
	if right == 0 {
		right = 1
	}
	t.moveCursor(t.buffer.cursorPos.Y, t.buffer.cursorPos.X+right)
}

func escapeMoveCursorLeft(t *Terminal, msg string) {
	left, _ := strconv.Atoi(msg)
	if left == 0 {
		left = 1
	}

	t.moveCursor(t.buffer.cursorPos.Y, t.buffer.cursorPos.X-left)
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
		// Clear from cursor
		row := t.buffer.Row(t.buffer.cursorPos.Y)
		from := t.buffer.savedCursorPos.X
		if t.buffer.cursorPos.X >= len(row.Chars) {
			from = len(row.Chars) - 1
		}
		if from > 0 {
			t.buffer.SetRow(t.buffer.cursorPos.Y, row.Chars[:from])
		} else {
			t.buffer.SetRow(t.buffer.cursorPos.Y, []Char{})
		}

		for i := t.buffer.cursorPos.Y + 1; i < len(t.buffer.lines); i++ {
			t.buffer.SetRow(i, []Char{})
		}
	case 1:
		// Clear to cursor
		row := t.buffer.Row(t.buffer.cursorPos.Y)
		chars := make([]Char, t.buffer.cursorPos.X)
		if t.buffer.cursorPos.X < len(row.Chars) {
			chars = append(chars, row.Chars[t.buffer.cursorPos.X:]...)
		}
		t.buffer.SetRow(t.buffer.cursorPos.Y, chars)

		for i := 0; i < t.buffer.cursorPos.Y-1; i++ {
			t.buffer.SetRow(i, []Char{})
		}
	case 2:
		t.Clear()
	}
}

func escapeInsertLines(t *Terminal, msg string) {
	rows, _ := strconv.Atoi(msg)
	if rows == 0 {
		rows = 1
	}
	i := t.scrollBottom
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

	t.scrollTop = start
	t.scrollBottom = end
}
