package jl

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// DefaultLogfmtPreferredFields is the set of fields that NewLogfmtPrinter orders ahead of other fields.
var DefaultLogfmtPreferredFields = []string{
	"timestamp",
	"time",
	"level",
	"thread",
	"logger",
	"message",
	"msg",
	"exceptions",
}

// LogfmtPrinter prints log entries in the logfmt format.
type LogfmtPrinter struct {
	// Out is the writer where formatted logs are written to.
	Out             io.Writer
	// PreferredFields is an order list of top-level keys that the logfmt formatter will display ahead of other
	// fields in the JSON log entry.
	PreferredFields []string
	// DisableColor disables ANSI color escape sequences.
	DisableColor    bool
}

// NewLogfmtPrinter allocates and returns a new LogFmtPrinter.
func NewLogfmtPrinter(w io.Writer) *LogfmtPrinter {
	return &LogfmtPrinter{
		Out:             w,
		PreferredFields: DefaultLogfmtPreferredFields,
	}
}

func (p *LogfmtPrinter) Print(input *Entry) {
	if input.Partials == nil {
		fmt.Fprintln(p.Out, string(input.Raw))
		return
	}
	entry := newLogfmtEntry(input, p.PreferredFields)
	color := entry.Color()

	sortedFields := append(entry.preferredFields, entry.sortedFields...)
	for i, field := range sortedFields {
		if i != 0 {
			fmt.Fprint(p.Out, " ")
		}
		key := field.Key
		if !p.DisableColor {
			key = ColorText(color, field.Key)
		}
		fmt.Fprintf(p.Out, "%s=%s", key, toString(field.Value))
	}
	fmt.Fprintln(p.Out)
}

func toString(v json.RawMessage) string {
	var str string
	if err := json.Unmarshal(v, &str); err != nil {
		return string(v)
	}
	return str
}

type logfmtEntry struct {
	rawMessage      []byte
	partials        map[string]json.RawMessage
	sortedFields    []*field
	preferredFields []*field
}

func newLogfmtEntry(m *Entry, preferredFields []string) *logfmtEntry {
	var preferredKeys = stringSet(preferredFields)
	var preferred, sorted []*field
	for _, k := range DefaultLogfmtPreferredFields {
		if v, ok := m.Partials[k]; ok {
			preferred = append(preferred, newField(k, v))
		}
	}
	var sortedKeys = sortKeys(m.Partials)
	for _, k := range sortedKeys {
		if _, ok := preferredKeys[k]; ok {
			continue
		}
		v := m.Partials[k]
		sorted = append(sorted, newField(k, v))
	}
	return &logfmtEntry{
		rawMessage:      m.Raw,
		partials:        m.Partials,
		sortedFields:    sorted,
		preferredFields: preferred,
	}
}

func (e *logfmtEntry) Color() Color {
	level := "info"
	if levelField, ok := e.partials["level"]; ok {
		level = toString(levelField)
	}
	if color, ok := LevelColors[strings.ToLower(level)]; ok {
		return color
	}
	return Green
}

type field struct {
	Key   string
	Value json.RawMessage
}

func newField(k string, v json.RawMessage) *field {
	return &field{Key: k, Value: v}
}

func stringSet(l []string) map[string]interface{} {
	m := make(map[string]interface{})
	for _, v := range l {
		m[v] = nil
	}
	return m
}

func sortKeys(m map[string]json.RawMessage) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.StringSlice(keys).Sort()
	return keys
}
