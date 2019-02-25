package jl

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type colorAssigner struct {
	assigned map[string]Color
	colorKey int
}

func newColorAssigner() *colorAssigner {
	return &colorAssigner{
		assigned: make(map[string]Color),
	}
}

func (a *colorAssigner) GetColor(v string) Color {
	if color, ok := a.assigned[v]; ok {
		return color
	}
	color := AllColors[a.colorKey%len(AllColors)]
	a.colorKey++
	a.assigned[v] = color
	return color
}

type CompactPrinter struct {
	w            io.Writer
	log          *log.Logger
	threadColors *colorAssigner
	loggerColors *colorAssigner
}

func NewCompactPrinter(w io.Writer) *CompactPrinter {
	return &CompactPrinter{
		w:            w,
		log:          log.New(w, "jl/formatter", log.LstdFlags),
		threadColors: newColorAssigner(),
		loggerColors: newColorAssigner(),
	}
}

func (h *CompactPrinter) Print(m *Line) {
	if m.JSON == nil {
		fmt.Fprintln(h.w, m.Raw)
		return
	}
	entry := newEntry(m, SpecialFields)
	h.printColored(entry)
}

func (h *CompactPrinter) printColored(entry *Entry) {
	levelColor := LevelColor(entry.Level())

	fmt.Fprint(h.w, ColorText(levelColor, strings.ToUpper(entry.Level().String())[0:4]))
	if timestamp, ok := entry.fieldMap["timestamp"]; ok {
		fmt.Fprintf(h.w, " %v", timestamp)
	}
	fmt.Fprint(h.w, " ")
	if thread, ok := entry.fieldMap["thread"]; ok {
		text := fmt.Sprint(thread)
		color := h.threadColors.GetColor(text)
		fmt.Fprintf(h.w, "[%v]\t", ColorText(color, text))
	}
	if logger, ok := entry.fieldMap["logger"]; ok {
		loggerStr := fmt.Sprint(logger)
		parts := strings.Split(loggerStr, ".")
		loggerStr = parts[len(parts)-1]
		color := h.threadColors.GetColor(loggerStr)
		fmt.Fprint(h.w, ColorText(color, loggerStr))
	}
	fmt.Fprint(h.w, "| ")
	if message, ok := entry.fieldMap["message"]; ok {
		fmt.Fprintf(h.w, "%v ", fmt.Sprint(message))
	}

	/*	for i, field := range entry.sortedFields {
			if i != 0 {
				fmt.Fprint(h.w, " ")
			}
			fmt.Fprintf(h.w, "%s=%s", ColorText(levelColor, field.Key), field.Value)
		}*/
	fmt.Fprintln(h.w)
}
