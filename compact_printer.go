package jl

import (
	"encoding/json"
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
		fmt.Fprintln(h.w, string(m.Raw))
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
	// End the log line
	fmt.Fprintln(h.w)

	// Exceptions go after the current log line
	if exceptions, ok := entry.fieldMap["exceptions"]; ok {
		var java struct {
			Exceptions []*KotlinException `json:"exceptions"`
		}
		if err := json.Unmarshal(entry.rawMessage, &java); err != nil {
			fmt.Println(err)
			fmt.Fprintln(h.w, "\t", exceptions)
		}
		for i, e := range java.Exceptions {
			fmt.Fprint(h.w, "  ")
			if i != 0 {
				fmt.Fprint(h.w, "Caused by: ")
			}
			fmt.Fprintf(h.w, "%s.%s: %s\n", e.Module, e.Type, e.Message)
			for _, stack := range e.StackTrace {
				fmt.Fprintf(h.w, "    at %s.%s(%s.%d)\n", stack.Module, stack.Func, stack.File, stack.Line)
			}
			if e.FramesOmitted > 0 {
				fmt.Fprintf(h.w, "    ...%d frames omitted...\n", e.FramesOmitted)
			}
		}
	}
}
