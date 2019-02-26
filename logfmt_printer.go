package jl

import (
	"fmt"
	"io"
	"log"
)

type LogfmtPrinter struct {
	w   io.Writer
	log *log.Logger
}

func NewLogfmtPrinter(w io.Writer) *LogfmtPrinter {
	return &LogfmtPrinter{
		w:   w,
		log: log.New(w, "jl/formatter", log.LstdFlags),
	}
}

func (h *LogfmtPrinter) Print(m *Line) {
	if m.JSON == nil {
		fmt.Fprintln(h.w, string(m.Raw))
		return
	}
	entry := newEntry(m, SpecialFields)
	h.printColored(entry)
}

func (h *LogfmtPrinter) printColored(entry *Entry) {
	levelColor := LevelColor(entry.Level())

	sortedFields := append(entry.specialFields, entry.sortedFields...)
	for i, field := range sortedFields {
		if i != 0 {
			fmt.Fprint(h.w, " ")
		}
		fmt.Fprintf(h.w, "%s=%s", ColorText(levelColor, field.Key), field.Value)
	}
	fmt.Fprintln(h.w)
}
