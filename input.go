package main

import "fyne.io/fyne/v2"

func (t *Terminal) onTypedKey(e *fyne.KeyEvent) {
	if e.Name == fyne.KeyEnter || e.Name == fyne.KeyReturn {
		_, _ = t.in.Write([]byte{'\r'})
	}
	if e.Name == fyne.KeyBackspace {
		_, _ = t.in.Write([]byte{asciiBackspace})
	}

}

func (t *Terminal) onTypedRune(r rune) {
	_, _ = t.in.Write([]byte{byte(r)})
}

func setInputs() {

}
