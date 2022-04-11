package term

import (
	"log"
	"strconv"
	"strings"

	"golang.org/x/image/colornames"
)

type parseState struct {
	parseMode bool

	// Control Sequence Introducer
	csi bool

	s string
}

func (buf *Buffer) Write(out []byte) {
	buf.Lock()
	defer buf.Unlock()
	state := parseState{}

	runes := []rune(string(out))
	char := Char{}
	for _, r := range runes {

		switch {
		case r == '\n':
			buf.insertLine(Line{})
		case r == 7:
			// bell
		case r == 8:
			index := len(buf.lines) - 1
			if index >= 0 {
				buf.lines[index].pop()
			}
		case r == 27:
			state.parseMode = true
		case r == '[' && state.parseMode:
			state.s = ""
			state.csi = true
		case r == 'm' && state.csi:
			for _, s := range strings.Split(state.s, ";") {
				switch {
				case s == "":
					continue
				case s == "0":
					char.BgColor, _ = parseBg(40)
					char.FgColor, _ = parseFg(37)
					continue
				case s == "1" || s == "01":
					continue
				case s == "39":
					char.FgColor, _ = parseFg(37)
					continue
				case s == "49":
					char.BgColor, _ = parseBg(37)
					continue
				default:
					i, err := strconv.Atoi(s)
					if err != nil {
						log.Println(err, "code:", s)
						continue
					}
					fgColor, ok := parseFg(i)
					if ok {
						char.FgColor = fgColor
						continue
					}
					bgColor, ok := parseBg(i)
					if ok {
						char.BgColor = bgColor
						continue
					}
					log.Println("ANSI code not implemented:", i)
				}
			}
			state.parseMode = false
			state.csi = false

		case state.csi || state.parseMode:
			state.s += string(r)
		case !state.csi && !state.parseMode:
			char.R = r

			if char.FgColor == nil {
				char.FgColor = colornames.White
			}

			lines := len(buf.lines)
			if lines == 0 {
				buf.insertLine(Line{})
			}
			buf.lines[len(buf.lines)-1].push(char)
		}
	}

}
