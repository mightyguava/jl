package jl

import "fmt"

type Color int

// Foreground text colors
const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// Foreground Hi-Intensity text colors
const (
	HiBlack Color = iota + 90
	HiRed
	HiGreen
	HiYellow
	HiBlue
	HiMagenta
	HiCyan
	HiWhite
)

var AllColors = []Color{
	// Skipping black because it's invisible on dark terminal backgrounds.
	Red,
	Green,
	Yellow,
	Blue,
	Magenta,
	Cyan,
	White,
	HiBlack,
	HiRed,
	HiGreen,
	HiYellow,
	HiBlue,
	HiMagenta,
	HiCyan,
	HiWhite,
}

func ColorText(c Color, text string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", c, text)
}

func LevelColor(level Level) Color {
	var levelColor Color
	switch level {
	case DebugLevel, TraceLevel:
		levelColor = White
	case WarnLevel:
		levelColor = Yellow
	case ErrorLevel, FatalLevel, PanicLevel:
		levelColor = Red
	default:
		levelColor = Green
	}
	return levelColor
}
