package jl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"unicode"
)

// FieldFmt specifies a single field formatted by the CompactPrinter.
type FieldFmt struct {
	// Name of the field. This is used to find the field by key name if Finders is not set.
	Name string
	// List of FieldFinders to use to locate the field. Finders are executed in order until the first one that returns
	// non-nil.
	Finders []FieldFinder
	// Takes the output of the Finder and turns it into a string. If not set, DefaultStringer is used.
	Stringer Stringer
	// List of transformers to run on the field found to format the field.
	Transformers []Transformer
}

// DefaultCompactPrinterFieldFmt is a format for the CompactPrinter that tries to present logs in an easily skimmable manner
// for most types of logs.
var DefaultCompactPrinterFieldFmt = []FieldFmt{{
	Name:         "level",
	Transformers: []Transformer{Truncate(4), UpperCase, ColorMap(LevelColors)},
}, {
	Finders: []FieldFinder{ByNames("timestamp", "time")},
}, {
	Name:         "thread",
	Transformers: []Transformer{Ellipsize(16), Format("[%s]"), RightPad(18), ColorSequence(AllColors)},
}, {
	Name:         "logger",
	Transformers: []Transformer{Ellipsize(20), Format("%s|"), LeftPad(21), ColorSequence(AllColors)},
}, {
	Finders: []FieldFinder{ByNames("message", "msg")},
}, {
	Finders:  []FieldFinder{JavaExceptionFinder, LogrusErrorFinder, ByNames("exceptions", "error")},
	Stringer: ErrorStringer,
}}

// CompactPrinter can print logs in a variety of compact formats, specified by FieldFormats.
type CompactPrinter struct {
	Out io.Writer
	// Disable colors disables adding color to fields.
	DisableColor bool
	// Disable truncate disables the Ellipsize and Truncate transforms.
	DisableTruncate bool
	// FieldFormats specifies the format the printer should use for logs. It defaults to DefaultCompactPrinterFieldFmt. Fields
	// are formatted in the order they are provided. If a FieldFmt produces a field that does not end with a whitespace,
	// a space character is automatically appended.
	FieldFormats []FieldFmt
}

// NewCompactPrinter allocates and returns a new compact printer.
func NewCompactPrinter(w io.Writer) *CompactPrinter {
	return &CompactPrinter{
		Out:            w,
		FieldFormats: DefaultCompactPrinterFieldFmt,
	}
}

func (p *CompactPrinter) Print(entry *Entry) {
	if entry.Partials == nil {
		fmt.Fprintln(p.Out, string(entry.Raw))
		return
	}
	for _, fieldFmt := range p.FieldFormats {
		ctx := Context{
			DisableColor:    p.DisableColor,
			DisableTruncate: p.DisableTruncate,
		}
		p.Out.Write([]byte(fieldFmt.format(&ctx, entry)))
	}
	p.Out.Write([]byte("\n"))
}

func (f *FieldFmt) format(ctx *Context, entry *Entry) string {
	var v interface{}
	// Find the value
	if len(f.Finders) > 0 {
		for _, finder := range f.Finders {
			if v = finder(entry); v != nil {
				break
			}
		}
	} else {
		v = entry.Partials[f.Name]
	}
	if v == nil {
		return ""
	}

	// Stringify the value
	var s string
	if f.Stringer != nil {
		s = f.Stringer(ctx, v)
	} else {
		s = DefaultStringer(ctx, v)
	}

	original := s
	ctx.Original = original
	// Apply transforms
	for _, transform := range f.Transformers {
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

// FieldFinder locates a field in the Entry and returns it.
type FieldFinder func(entry *Entry) interface{}

// ByNames locates fields by their top-level key name in the JSON log entry, and returns the field as a json.RawMessage.
func ByNames(names ...string) FieldFinder {
	return func(entry *Entry) interface{} {
		for _, name := range names {
			if v, ok := entry.Partials[name]; ok {
				return v
			}
		}
		return nil
	}
}

// LogrusErrorFinder finds logrus error in the JSON log and returns it as a LogrusError.
func LogrusErrorFinder(entry *Entry) interface{} {
	var errStr, stack string
	if errV, ok := entry.Partials["error"]; !ok {
		return nil
	} else if err := json.Unmarshal(errV, &errStr); err != nil {
		return nil
	}
	if stackV, ok := entry.Partials["stack"]; !ok {
		return nil
	} else if err := json.Unmarshal(stackV, &stack); err != nil {
		return nil
	}
	return LogrusError{errStr, stack}
}

// JavaExceptionFinder finds a Java exception containing a stracetrace and returns it as a JavaExceptions.
func JavaExceptionFinder(entry *Entry) interface{} {
	var java struct {
		Exceptions []*JavaException `json:"exceptions"`
	}
	if err := json.Unmarshal(entry.Raw, &java); err == nil {
		return JavaExceptions(java.Exceptions)
	}
	return nil
}

// Stringer transforms a field returned by the FieldFinder into a string.
type Stringer func(ctx *Context, v interface{}) string

var _ = Stringer(DefaultStringer)
var _ = Stringer(ErrorStringer)

// DefaultStringer attempts to turn a field into string by attempting the following in order
// 1. casting it to a string
// 2. unmarshalling it as a json.RawMessage
// 3. using fmt.Sprintf("%v", input)
func DefaultStringer(ctx *Context, v interface{}) string {
	var s string
	if tmp, ok := v.(string); ok {
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
	return s
}

// ErrorStringer stringifies LogrusError, JavaExceptions to a multiline string. If the field is neither, it falls back
// to the DefaultStringer.
func ErrorStringer(ctx *Context, v interface{}) string {
	w := &bytes.Buffer{}
	if logrusErr, ok := v.(LogrusError); ok {
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
				fmt.Fprintf(w, "\n    at %s.%s(%s:%d)", stack.Module, stack.Func, stack.File, stack.Line)
			}
			if e.FramesOmitted > 0 {
				fmt.Fprintf(w, "\n    ...%d frames omitted...", e.FramesOmitted)
			}
		}
		return w.String()
	} else {
		return DefaultStringer(ctx, v)
	}
}

// Context provides the current transformation context, to be used by Transformers and Stringers.
type Context struct {
	// The original string before any transformations were applied.
	Original string
	// Indicates that terminal color escape sequences should be disabled.
	DisableColor bool
	// Indicates that fields should not be truncated.
	DisableTruncate bool
}
