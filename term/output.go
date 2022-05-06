package term

const (
	asciBell       = 7
	asciiBackspace = 8
	asciiEscape    = 27
	asciiCarriage  = '\r'
	asciiNewLine   = '\n'
	asciiTab       = '\t'
	noEscape = -100
	tabWidth = 8
)

type parseState struct {
	vt100 rune

	esc int

	osc bool

	s string
}

func (t *Terminal) handleOutput(out []byte) {
	var previous *parseState
	state := &parseState{}
	if previous != nil {
		state = previous
		previous = nil
	} else {
		state.esc = noEscape
	}

	runes := []rune(string(out))
	for i, r := range runes {
		if r == asciiEscape {
			state.esc = i
			continue
		}
		if state.esc == i-1 {
			if r == '[' {
				continue
			}

			switch r {
			case '\\':
				t.handleOSC(state.s)
				state.s = ""
				state.osc = false
			case ']':
				state.osc = true
			case '(', ')':
				state.vt100 = r
			case '7':
				t.buffer.savedCursorPos.X = t.buffer.cursorPos.X
				t.buffer.savedCursorPos.Y = t.buffer.cursorPos.Y
			case '8':
				t.buffer.savedCursorPos.X = t.buffer.cursorPos.X
				t.buffer.savedCursorPos.Y = t.buffer.cursorPos.Y
			case 'D':
				t.ScrollDown()
			case 'M':
				t.ScrollUp()
			case '=', '>':
			}
			state.esc = noEscape
			continue
		}
		if state.osc {
			if r == asciBell || r == 0 {
				t.handleOSC(state.s)
				state.s = ""
				state.osc = false
			} else {
				state.s += string(r)
			}
			continue
		}
		if state.vt100 != 0 {
			state.vt100 = 0
			continue
		}
		if state.esc != noEscape {
			state.s += string(r)
			if (r < '0' || r > '9') && r != ';' && r != '=' && r != '?' {
				t.handleEscape(state.s)
				state.s = ""
				state.esc = noEscape
			}
			continue
		}

		switch {
		case r == asciiCarriage:
			t.moveCursor(t.buffer.cursorPos.Y, 0)

		case r == asciiBackspace:
			t.Backspace()

		case r == asciBell:

		case r == 0x0e || r == 0x0f:
			continue

		case r == asciiNewLine:
			t.moveCursor(t.buffer.cursorPos.Y+1, t.buffer.cursorPos.X)

		case r == asciiTab:
			end := t.buffer.cursorPos.X + tabWidth

			for t.buffer.cursorPos.X < end {

				t.buffer.insertChar(Char{
					R:       ' ',
					FgColor: t.currentFG,
					BgColor: t.currentBG,
				})
			}

		default:
			t.buffer.insertChar(Char{
				R:       r,
				FgColor: t.currentFG,
				BgColor: t.currentBG,
			})
		}
	}

	if state.esc != noEscape {
		state.esc = -1 - (len(state.s))
		previous = state
	}
}
