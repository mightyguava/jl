package jl

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"
)

var (
	UpperCase = TransformFunc(strings.ToUpper)
	LowerCase = TransformFunc(strings.ToLower)
)

type Transformer interface {
	Transform(ctx *Context, acc string) string
}

type TransformFunc func(string) string

func (f TransformFunc) Transform(ctx *Context, acc string) string {
	return f(acc)
}

type Truncate int

func (t Truncate) Transform(ctx *Context, acc string) string {
	if utf8.RuneCountInString(acc) <= int(t) {
		return acc
	}
	return acc[:t]
}

type Ellipsize int

func (remain Ellipsize) Transform(ctx *Context, acc string) string {
	length := utf8.RuneCountInString(acc)
	if length <= int(remain) {
		return acc
	}
	remain -= 1 // account for the ellipsis
	chomped := length - int(remain)
	start := int(remain)/2
	end := start + chomped
	return acc[:start] + "…" + acc[end:]
}

type LeftPad int

func (t LeftPad) Transform(ctx *Context, acc string) string {
	spaces := int(t) - utf8.RuneCountInString(acc)
	if spaces <= 0 {
		return acc
	}
	buf := bytes.NewBuffer(make([]byte, spaces+len(acc)))
	for i := 0; i < spaces; i++ {
		buf.WriteRune(' ')
	}
	buf.WriteString(acc)
	return buf.String()
}

type RightPad int

func (t RightPad) Transform(ctx *Context, acc string) string {
	pad := int(t) - utf8.RuneCountInString(acc)
	if pad <= 0 {
		return acc
	}
	buf := bytes.NewBuffer(make([]byte, pad+len(acc)))
	buf.WriteString(acc)
	for i := 0; i < pad; i++ {
		buf.WriteRune(' ')
	}
	return buf.String()
}

type Format string

func (t Format) Transform(ctx *Context, acc string) string {
	return fmt.Sprintf(string(t), acc)
}
