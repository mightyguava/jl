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
		var data map[string]interface{}
		_ = json.Unmarshal(raw, &data)
		message := &Line{
			JSON: data,
			Raw:  raw,
		}
		p.h.Print(message)
	}
	return p.scan.Err()
}

type EntryPrinter interface {
	Print(*Line)
}

type Line struct {
	JSON map[string]interface{}
	Raw  []byte
}