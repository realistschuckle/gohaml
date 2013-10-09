package parser

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type NodeVisitor interface {
	VisitDoctype(*DoctypeNode)
	VisitTag(*TagNode)
}

type Node interface {
	Accept(NodeVisitor)
}

type HamlParser interface {
	Parse(io.RuneReader) (ParsedDoc, ParseError)
}

type LineParser interface {
	Parse([]rune) (Node, *ParseError)
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
	var parser LineParser

	for r, _, ok := scanner.ReadRune(); ok == nil; r, _, ok = scanner.ReadRune() {
		line = append(line, r)
		if r == '\n' {
			if line[0] == '!' {
				parser = &DoctypeParser{}
			} else {
				parser = &TagParser{}
			}
			n, e := parser.Parse(line)
			if e != nil {
				err = e
				return
			}
			nodes = append(nodes, n)
			line = linebuf[0:0]
		}
	}
	if len(line) > 0 {
		if line[0] == '!' {
			parser = &DoctypeParser{}
		} else {
			parser = &TagParser{}
		}
		n, e := parser.Parse(line)
		if e != nil {
			err = e
			return
		}
		nodes = append(nodes, n)
	}

	doc = ParsedDoc{nodes}
	return
}

type DoctypeParser struct {
}

func (self *DoctypeParser) Parse(input []rune) (n Node, err *ParseError) {
	if len(input) < 3 || input[0] != '!' || input[1] != '!' || input[2] != '!' {
		err = &ParseError{1, 1}
		return
	}
	specifier := strings.TrimFunc(string(input[3:]), func(r rune) bool {
		return unicode.IsSpace(r)
	})
	n = &DoctypeNode{specifier}
	return
}

type TagParser struct {
}

func (self *TagParser) Parse(input []rune) (n Node, err *ParseError) {
	tn := &TagNode{"div", "", nil, nil, nil, false}
	if input[0] != '%' {
		err = &ParseError{1, 1}
		return
	}

	start := 0
	for i := 0; i < len(input); i += 1 {
		if input[i] == '%' {
			start = i + 1
			for i = i + 1; i < len(input); i += 1 {
				if !unicode.IsLetter(input[i]) &&
				   !unicode.IsDigit(input[i]) &&
				   input[i] != '-' &&
				   input[i] != '_' &&
				   input[i] != ':' {
				   	tn.Name = string(input[start:i])
				   	break
			   }
			   if i == len(input) - 1 {
				   	tn.Name = string(input[start:i + 1])
			   }
			}
			i -= 1
			continue
		}
		if input[i] == '.' {
			start = i + 1
			for i = i + 1; i < len(input); i += 1 {
				if !unicode.IsLetter(input[i]) &&
				   !unicode.IsDigit(input[i]) &&
				   input[i] != '-' &&
				   input[i] != '_' {
				   	class := string(input[start:i])
				   	tn.Classes = append(tn.Classes, class)
				   	break
			   }
			   if i == len(input) - 1 {
				   	class := string(input[start:i + 1])
				   	tn.Classes = append(tn.Classes, class)
			   }
			}
			i -= 1
			continue
		}
		if input[i] == '#' {
			start = i + 1
			for i = i + 1; i < len(input); i += 1 {
				if !unicode.IsLetter(input[i]) &&
				   !unicode.IsDigit(input[i]) &&
				   input[i] != '-' &&
				   input[i] != '_' {
				   	tn.Id = string(input[start:i])
				   	break
				}
				if i == len(input) - 1 {
					tn.Id = string(input[start:i + 1])
				}
			}
			i -= 1
			continue
		}
		if input[i] == '/' {
			tn.Close = true
		}
	}

	n = tn
	return
}

type DoctypeNode struct {
	Specifier string
}

func (self *DoctypeNode) Accept(visitor NodeVisitor) {
	visitor.VisitDoctype(self)
}

type TagNode struct {
	Name     string
	Id       string
	Classes  []string
	Attrs    map[string]string
	Children []Node
	Close    bool
}

func (self *TagNode) Accept(visitor NodeVisitor) {
	visitor.VisitTag(self)
}
