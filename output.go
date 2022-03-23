package main

func (t *Terminal) handleOutput(out []byte) {
	for _, r := range []rune(string(out)) {

		if t.lines >= len(t.content) {
			t.content = append(t.content, []rune{})
		}
		if r == asciiBackspace {
			content := t.content[t.lines]
			t.content[t.lines] = content[:len(content)-1]
			continue
		}
		if r == asciiBell {
			continue
		}

		if r == '\n' {
			t.lines++
			continue
		}

		t.content[t.lines] = append(t.content[t.lines], r)
	}
}
