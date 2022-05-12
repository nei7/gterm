package term

func (t *Terminal) handleOSC(code string) {
	if len(code) <= 2 || code[1] != ';' {
		return
	}

	switch code[0] {
	case '0':
		t.setTitle(code[2:])
	case '1':
	case '2':
		t.setTitle(code[2:])
	case '7':
	}
}
