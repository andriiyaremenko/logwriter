package color

import "regexp"

type Color string

const (
	ANSIReset       Color = "\033[0m"
	ANSIColorRed    Color = "\033[31m"
	ANSIColorGreen  Color = "\033[32m"
	ANSIColorYellow Color = "\033[33m"
	ANSIColorBlue   Color = "\033[34m"
	ANSIColorPurple Color = "\033[35m"
	ANSIColorCyan   Color = "\033[36m"
	ANSIColorWhite  Color = "\033[37m"
	ANSIColorGray   Color = "\033[90m"
	ANSIFontBold    Color = "\033[1m"

	ColorTrace Color = ANSIColorGray
	ColorDebug Color = ANSIColorCyan
	ColorInfo  Color = ANSIColorGreen
	ColorWarn  Color = ANSIColorYellow
	ColorError Color = ANSIColorRed
	ColorFatal Color = ANSIFontBold + ANSIColorRed
)

var ansiiColorMatch = regexp.MustCompile("\u001B\\[[;\\d]*m")

// Returns text prepended by ANSI color code and appended by ANSI color reset code.
func ColorizeText(color Color, text string) string {
	return string(color) + text + string(ANSIReset)
}

// Returns text without color.
func ClearColors(text string) string {
	return ansiiColorMatch.ReplaceAllString(text, "")
}

func GetLevelColor(level int) Color {
	switch level {
	case 1:
		return ColorDebug
	case 2:
		return ColorInfo
	case 3:
		return ColorWarn
	case 4:
		return ColorError
	}

	if level < 1 {
		return ColorTrace
	}

	return ColorFatal
}
