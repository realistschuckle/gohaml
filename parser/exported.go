package parser

import (
	"io"
)

type ParsedDoc struct {
	// Indentation
}

type HamlParser interface {
	Parse(io.RuneReader) (ParsedDoc, error)
}

type DefaultParser struct {
}

func (self *DefaultParser) Parse(input io.RuneReader) (doc ParsedDoc, err error) {
	scanner := scanner{input, [8]rune{}, 0, 0}
	linebuf := [1000]rune{}
	line := linebuf[0:0]
	
	for r, _, ok := scanner.ReadRune(); ok == nil; r, _, ok = scanner.ReadRune() {
		line = append(line, r)
	}

	doc = ParsedDoc{}
	return
}
