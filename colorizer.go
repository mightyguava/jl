package jl

import (
	"strings"
)

type sequentialColorizer struct {
	assigned map[string]Color
	seq      int
	colors   []Color
}

func ColorSequence(colors []Color) *sequentialColorizer {
	return &sequentialColorizer{
		assigned: make(map[string]Color),
		colors:   colors,
	}
}

func (a *sequentialColorizer) Transform(ctx *Context, acc string) string {
	if color, ok := a.assigned[ctx.Original]; ok {
		return ColorText(color, acc)
	}
	color := a.colors[a.seq%len(AllColors)]
	a.seq++
	a.assigned[ctx.Original] = color
	return ColorText(color, acc)
}

type mappingColorizer struct {
	mapping map[string]Color
}

func ColorMap(mapping map[string]Color) *mappingColorizer {
	lowered := make(map[string]Color, len(mapping))
	for k, v := range mapping {
		lowered[strings.ToLower(k)] = v
	}
	return &mappingColorizer{lowered}
}

func (c *mappingColorizer) Transform(ctx *Context, acc string) string {
	if color, ok := c.mapping[strings.ToLower(ctx.Original)]; ok {
		return ColorText(color, acc)
	}
	return acc
}
