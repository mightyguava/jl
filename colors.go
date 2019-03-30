package jl

import (
	"fmt"
)

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

// AllColors is the set of colors used by default by DefaultCompactFieldFmts for ColorSequence.
var AllColors = []Color{
	// Skipping black because it's invisible on dark terminal backgrounds.
	// Skipping red because it's too prominent and means error
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

// ColorText wraps a text with ANSI escape codes to produce terminal colors.
func ColorText(c Color, text string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", c, text)
}

// LevelColors is a mapping of log level strings to colors.
var LevelColors = map[string]Color{
	"trace": White,
	"debug": White,
	"info": Green,
	"warn": Yellow,
	"warning": Yellow,
	"error": Red,
	"fatal": Red,
	"panic": Red,
}
