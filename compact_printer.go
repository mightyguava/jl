package jl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"unicode"
)

var CompactPrinterFieldFmt = []FieldFmt{{
	Name:       "level",
	Transforms: []Transformer{Truncate(4), UpperCase, ColorMap(LevelColors)},
}, {
	Finders: []FieldFinder{ByNames("timestamp", "time")},
}, {
	Name:       "thread",
	Transforms: []Transformer{Ellipsize(16), Format("[%s]"), RightPad(18), ColorSequence(AllColors)},
}, {
	Name:       "logger",
	Transforms: []Transformer{Ellipsize(20), Format("%s|"), LeftPad(21), ColorSequence(AllColors)},
}, {
	Name: "message",
}, {
	Name:     "error",
	Finders:  []FieldFinder{JavaExceptionFinder, LogrusErrorFinder, ByNames("exceptions", "error")},
	Stringer: ErrorStringer,
}}

type CompactPrinter struct {
	w io.Writer
	// Log is used for logging parsing and rendering errors
	Log *log.Logger
	// Disable colors disables adding color to fields
	DisableColors bool
	FieldFormats  []FieldFmt
}

func NewCompactPrinter(w io.Writer) *CompactPrinter {
	return &CompactPrinter{
		w:            w,
		Log:          log.New(w, "jl/formatter", log.LstdFlags),
		FieldFormats: CompactPrinterFieldFmt,
	}
}

func (p *CompactPrinter) Print(m *Line) {
	if m.Partials == nil {
		fmt.Fprintln(p.w, string(m.Raw))
		return
	}
	entry := newEntry(m, SpecialFields)
	for _, fieldFmt := range p.FieldFormats {
		p.w.Write([]byte(fieldFmt.format(entry)))
	}
	p.w.Write([]byte("\n"))
}

type FieldFmt struct {
	Name       string
	Transforms []Transformer
	Finders             []FieldFinder
	Stringer            Stringer
}

func (f *FieldFmt) format(entry *Entry) string {
	var v interface{}
	// Find the value
	if len(f.Finders) > 0 {
		for _, finder := range f.Finders {
			if v = finder(entry); v != nil {
				break
			}
		}
	} else {
		v = entry.partials[f.Name]
	}
	if v == nil {
		return ""
	}

	ctx := &Context{}
	// Stringify the value
	var s string
	if f.Stringer != nil {
		s = f.Stringer(ctx, v)
	} else if tmp, ok := v.(string); ok {
		s = tmp
	} else if rawMsg, ok := v.(json.RawMessage); ok {
		var unmarshaled interface{}
		if err := json.Unmarshal(rawMsg, &unmarshaled); err != nil {
			s = string(rawMsg)
		} else {
			s = fmt.Sprintf("%v", unmarshaled)
		}
	} else {
		s = fmt.Sprintf("%v", v)
	}

	original := s
	ctx.Original = original
	// Apply transforms
	for _, transform := range f.Transforms {
		s = transform.Transform(ctx, s)
	}

	if s == "" {
		return ""
	}

	// Add a space if needed
	if !unicode.IsSpace(rune(s[len(s)-1])) {
		s += " "
	}

	return s
}

type FieldFinder func(entry *Entry) interface{}

func ByNames(names ...string) FieldFinder {
	return func(entry *Entry) interface{} {
		for _, name := range names {
			if v, ok := entry.partials[name]; ok {
				return v
			}
		}
		return nil
	}
}

func LogrusErrorFinder(entry *Entry) interface{} {
	var errStr, stack string
	if errV, ok := entry.partials["error"]; !ok {
		return nil
	} else if err := json.Unmarshal(errV, &errStr); err != nil {
		return nil
	}
	if stackV, ok := entry.partials["stack"]; !ok {
		return nil
	} else if err := json.Unmarshal(stackV, &stack); err != nil {
		return nil
	}
	return logrusError{errStr, stack}
}

func JavaExceptionFinder(entry *Entry) interface{} {
	var java struct {
		Exceptions []*JavaException `json:"exceptions"`
	}
	if err := json.Unmarshal(entry.rawMessage, &java); err == nil {
		return JavaExceptions(java.Exceptions)
	}
	return nil
}

type Stringer func(ctx *Context, v interface{}) string

func ErrorStringer(ctx *Context, v interface{}) string {
	w := &bytes.Buffer{}
	if logrusErr, ok := v.(logrusError); ok {
		// left pad with a tab
		lines := strings.Split(logrusErr.Error, "\n")
		stackStr := "\t" + strings.Join(lines, "\n\t")
		w.WriteString(stackStr)
		return w.String()
	} else if exceptions, ok := v.(JavaExceptions); ok {
		for i, e := range []*JavaException(exceptions) {
			fmt.Fprint(w, "\n  ")
			if i != 0 {
				fmt.Fprint(w, "Caused by: ")
			}
			msg := e.Message
			if !ctx.DisableColor {
				msg = ColorText(Red, msg)
			}
			fmt.Fprintf(w, "%s.%s: %s", e.Module, e.Type, msg)
			for _, stack := range e.StackTrace {
				fmt.Fprintf(w, "\n    at %s.%s(%s.%d)", stack.Module, stack.Func, stack.File, stack.Line)
			}
			if e.FramesOmitted > 0 {
				fmt.Fprintf(w, "\n    ...%d frames omitted...", e.FramesOmitted)
			}
		}
		return w.String()
	} else if raw, ok := v.(json.RawMessage); ok {
		return string(raw)
	} else {
		return fmt.Sprintf("%v", raw)
	}
}

type logrusError struct {
	Error string
	Stack string
}

type Context struct {
	Original string
	DisableColor bool
}
