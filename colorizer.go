package jl

import (
	"strings"
)

type sequentialColorizer struct {
	assigned map[string]Color
	seq      int
	colors   []Color
}

// ColorSequence assigns colors to inputs sequentially. Once an input is seen and assigned a color, future iputs with
// the same value will always be assigned the same color.
func ColorSequence(colors []Color) *sequentialColorizer {
	return &sequentialColorizer{
		assigned: make(map[string]Color),
		colors:   colors,
	}
}

func (a *sequentialColorizer) Transform(ctx *Context, input string) string {
	if ctx.DisableColor {
		return input
	}
	if color, ok := a.assigned[ctx.Original]; ok {
		return ColorText(color, input)
	}
	color := a.colors[a.seq%len(AllColors)]
	a.seq++
	a.assigned[ctx.Original] = color
	return ColorText(color, input)
}

type mappingColorizer struct {
	mapping map[string]Color
}

// ColorMap assigns colors by mapping the original, pre-transform field value to a color based on a pre-defined mapping.
func ColorMap(mapping map[string]Color) *mappingColorizer {
	lowered := make(map[string]Color, len(mapping))
	for k, v := range mapping {
		lowered[strings.ToLower(k)] = v
	}
	return &mappingColorizer{lowered}
}

func (c *mappingColorizer) Transform(ctx *Context, input string) string {
	if ctx.DisableColor {
		return input
	}
	if color, ok := c.mapping[strings.ToLower(ctx.Original)]; ok {
		return ColorText(color, input)
	}
	return input
}
