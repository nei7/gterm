package term

import (
	"image/color"
)

func parseBg(i int) (color.RGBA, bool) {
	bg := make(map[int]color.RGBA, 16)

	bg[40] = color.RGBA{0, 0, 0, 255}        // Black
	bg[41] = color.RGBA{255, 0, 0, 255}      // Red
	bg[42] = color.RGBA{0, 255, 0, 255}      // Green
	bg[43] = color.RGBA{170, 85, 0, 255}     // Yellow
	bg[44] = color.RGBA{0, 0, 170, 255}      // Blue
	bg[45] = color.RGBA{170, 0, 170, 255}    // Magenta
	bg[46] = color.RGBA{0, 170, 170, 255}    // Cyan
	bg[47] = color.RGBA{170, 170, 170, 255}  // White
	bg[100] = color.RGBA{85, 85, 85, 255}    // Bright Black (Gray)
	bg[101] = color.RGBA{255, 85, 85, 255}   // Bright Red
	bg[102] = color.RGBA{85, 255, 85, 255}   // Bright Green
	bg[103] = color.RGBA{255, 255, 85, 255}  // Bright Yellow
	bg[104] = color.RGBA{85, 85, 255, 255}   // Bright Blue
	bg[105] = color.RGBA{255, 85, 255, 255}  // Bright Magenta
	bg[106] = color.RGBA{85, 255, 255, 255}  // Bright Cyan
	bg[107] = color.RGBA{255, 255, 255, 255} // Bright White

	c, ok := bg[i]
	return c, ok
}
func parseFg(i int) (color.RGBA, bool) {
	fg := make(map[int]color.RGBA, 16)

	fg[30] = color.RGBA{0, 0, 0, 255}       // Black
	fg[31] = color.RGBA{170, 0, 0, 255}     // Red
	fg[32] = color.RGBA{0, 170, 0, 255}     // Green
	fg[33] = color.RGBA{170, 85, 0, 255}    // Yellow
	fg[34] = color.RGBA{0, 0, 170, 255}     // Blue
	fg[35] = color.RGBA{170, 0, 170, 255}   // Magenta
	fg[36] = color.RGBA{0, 170, 170, 255}   // Cyan
	fg[37] = color.RGBA{170, 170, 170, 255} // White
	fg[90] = color.RGBA{85, 85, 85, 255}    // Bright Black (Gray)
	fg[91] = color.RGBA{255, 85, 85, 255}   // Bright Red
	fg[92] = color.RGBA{85, 255, 85, 255}   // Bright Green
	fg[93] = color.RGBA{255, 255, 85, 255}  // Bright Yellow
	fg[94] = color.RGBA{85, 85, 255, 255}   // Bright Blue
	fg[95] = color.RGBA{255, 85, 255, 255}  // Bright Magenta
	fg[96] = color.RGBA{85, 255, 255, 255}  // Bright Cyan
	fg[97] = color.RGBA{255, 255, 255, 255} // Bright White

	c, ok := fg[i]
	return c, ok
}
