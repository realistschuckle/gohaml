package parser

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type NodeVisitor interface {
	VisitDoctype(*DoctypeNode)
}

type Node interface {
	Accept(NodeVisitor)
}

type HamlParser interface {
	Parse(io.RuneReader) (ParsedDoc, ParseError)
}

type LineParser interface {
	Parse([]rune) (Node, ParseError)
}

type ParseError struct {
	line   int
	column int
}

func (self *ParseError) Error() (msg string) {
	msg = fmt.Sprintf("Error at (%i, %i)", self.line, self.column)
	return
}

type ParsedDoc struct {
	// Indentation
	Nodes []Node
}

func (self *ParsedDoc) Accept(visitor NodeVisitor) {
	for i := 0; i < len(self.Nodes); i += 1 {
		self.Nodes[i].Accept(visitor)
	}
	return
}

type DefaultParser struct {
}

func (self *DefaultParser) Parse(input io.RuneReader) (doc ParsedDoc, err error) {
	scanner := scanner{input, [8]rune{}, 0, 0}
	linebuf := [1000]rune{}
	line := linebuf[0:0]
	nodes := []Node{}

	for r, _, ok := scanner.ReadRune(); ok == nil; r, _, ok = scanner.ReadRune() {
		line = append(line, r)
		if r == '\n' {
			parser := DoctypeParser{}
			n, _ := parser.Parse(line)
			nodes = append(nodes, n)
			line = linebuf[0:0]
		}
	}
	if len(line) > 0 {
		parser := DoctypeParser{}
		n, _ := parser.Parse(line)
		nodes = append(nodes, n)
	}

	doc = ParsedDoc{nodes}
	return
}

type DoctypeParser struct {
}

func (self *DoctypeParser) Parse(input []rune) (n Node, err error) {
	if len(input) < 3 || input[0] != '!' || input[1] != '!' || input[2] != '!' {
		err = &ParseError{1, 1}
	}
	specifier := strings.TrimFunc(string(input[3:]), func(r rune) bool {
		return unicode.IsSpace(r)
	})
	n = &DoctypeNode{specifier}
	return
}

type DoctypeNode struct {
	Specifier string
}

func (self *DoctypeNode) Accept(visitor NodeVisitor) {
	visitor.VisitDoctype(self)
}
