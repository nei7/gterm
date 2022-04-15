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

var previous *parseState

func (t *Terminal) Print(out []byte) {
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
			t.moveCursor(t.Buffer.cursorPos.Y, 0)

		case r == asciiBackspace:
			t.Backspace()

		case r == asciBell:

		case r == asciiNewLine:
			t.Buffer.insertLine(Line{})
			t.moveCursor(t.Buffer.cursorPos.Y+1, t.Buffer.cursorPos.X)

		case r == asciiTab:
			end := t.Buffer.cursorPos.X - t.Buffer.cursorPos.X%tabWidth + tabWidth

			for t.Buffer.cursorPos.X < end {

				t.Buffer.appendToLine(t.Buffer.cursorPos.Y, Char{
					R:       ' ',
					FgColor: t.currentFG,
				})
			}
		default:
			t.Buffer.appendToLine(t.Buffer.cursorPos.Y, Char{
				R:       r,
				FgColor: t.currentFG,
			})

		}

	}

	if state.esc != noEscape {
		state.esc = -1 - (len(state.s))
		previous = state
	}
}
