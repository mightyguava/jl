package jl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

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

// ErrorStringer stringifies LogrusError to a multiline string. If the field is not a LogrusError, it falls back
// to the DefaultStringer.
func ErrorStringer(ctx *Context, v interface{}) string {
	w := &bytes.Buffer{}
	if logrusErr, ok := v.(LogrusError); ok {
		w.WriteString("\n  ")
		w.WriteString(logrusErr.Error)
		w.WriteRune('\n')
		// left pad with a tab
		lines := strings.Split(logrusErr.Stack, "\n")
		stackStr := "\t" + strings.Join(lines, "\n\t")
		w.WriteString(stackStr)
		return w.String()
	}  else {
		return DefaultStringer(ctx, v)
	}
}
