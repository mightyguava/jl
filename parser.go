package jl

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
)

type Parser struct {
	r    io.Reader
	scan *bufio.Scanner
	h    EntryPrinter
	log  *log.Logger
}

func NewParser(r io.Reader, h EntryPrinter) *Parser {
	return &Parser{
		r:    r,
		scan: bufio.NewScanner(r),
		h:    h,
		log:  log.New(os.Stderr, "jl/parser", log.LstdFlags),
	}
}

func (p *Parser) Consume() error {
	s := p.scan
	for s.Scan() {
		raw := s.Bytes()
		var partials map[string]json.RawMessage
		_ = json.Unmarshal(raw, &partials)
		message := &Line{
			Partials:    partials,
			Raw:         raw,
		}
		p.h.Print(message)
	}
	return p.scan.Err()
}

type EntryPrinter interface {
	Print(*Line)
}

type Line struct {
	Partials    map[string]json.RawMessage
	Raw         []byte
}
