package term

import (
	"log"
)

func (t *Terminal) handleOSC(code string) {

	if len(code) <= 2 || code[1] != ';' {
		return
	}

	switch code[0] {
	case '0':
		// set icon name
		t.setTitle(code[2:])
	case '1':
		// set icon name
	case '2':
		t.setTitle(code[2:])
	case '7':
		// set directory
	default:
		if t.debug {
			log.Println("Unrecognised OSC:", code)
		}
	}
}
