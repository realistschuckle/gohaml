package gohaml

import (
	"container/list"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type NodeVisitor interface {
	VisitDoctype(*DoctypeNode)
	VisitTag(*TagNode)
	VisitStaticLine(*StaticLineNode)
}

type Node interface {
	Accept(NodeVisitor)
	AddChild(Node) (ok bool)
}

type HamlParser interface {
	Parse(io.RuneReader) (ParsedDoc, ParseError)
	Indentation() string
}

type LineParser interface {
	Parse(string, []rune) (Node, bool, *ParseError)
}

type Attribute struct {
	Name  string
	Value Node
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
	indentation string
}

func (self *DefaultParser) Indentation() (s string) {
	s = self.indentation
	return
}

func (self *DefaultParser) Parse(input io.RuneReader) (doc ParsedDoc, err error) {
	scanner := scanner{input, [8]rune{}, 0, 0}
	linebuf := [1000]rune{}
	line := linebuf[0:0]
	lineNumber := 0
	nodes := []Node{}
	stack := list.New()
	indentDepth := 0
	var parser LineParser

	parseLine := func(line []rune) (e *ParseError) {
		var space string
		for i := 0; i < len(line); i += 1 {
			if !unicode.IsSpace(line[i]) {
				space = string(line[0:i])
				line = line[i:]
				break
			}
		}
		if len(self.indentation) == 0 && len(space) > 0 {
			self.indentation = space
		}
		if len(self.indentation) > 0 {
			indentDepth = len(space) / len(self.indentation)
		}
		if parser == nil {
			switch line[0] {
			case '!':
				parser = &DoctypeParser{}
			case '#', '.', '%':
				parser = &TagParser{}
			default:
				parser = &StaticParser{}
			}
		} else {
			space = ""
		}
		n, complete, e := parser.Parse(space, line)
		if !complete {
			return
		}
		for stack.Len() > indentDepth {
			stack.Remove(stack.Back())
		}
		if stack.Len() > 0 {
			parent := stack.Back().Value.(Node)
			parent.AddChild(n)
		} else {
			nodes = append(nodes, n)
		}
		stack.PushBack(n)
		if e != nil {
			e.column += len(space)
		}
		if complete {
			parser = nil
		}
		return
	}

	for r, _, ok := scanner.ReadRune(); ok == nil; r, _, ok = scanner.ReadRune() {
		line = append(line, r)
		if r == '\n' {
			lineNumber += 1
			e := parseLine(line)
			if e != nil {
				e.line = lineNumber
				err = e
				return
			}
			line = linebuf[0:0]
		}
	}
	if len(line) > 0 {
		e := parseLine(line)
		if e != nil {
			e.line = lineNumber
			err = e
			return
		}
	}

	doc = ParsedDoc{nodes}
	return
}

type DoctypeParser struct {
}

func (self *DoctypeParser) Parse(indent string, input []rune) (n Node, completed bool, err *ParseError) {
	if len(input) < 3 || input[0] != '!' || input[1] != '!' || input[2] != '!' {
		err = &ParseError{1, 1}
		return
	}
	specifier := strings.TrimFunc(string(input[3:]), func(r rune) bool {
		return unicode.IsSpace(r)
	})
	n = &DoctypeNode{specifier}
	completed = true
	return
}

type TagParser struct {
	tag       *TagNode
	completed bool
}

func (self *TagParser) Parse(indent string, input []rune) (n Node, completed bool, err *ParseError) {
	if self.tag == nil {
		self.completed = true
		self.tag = &TagNode{"div", "", nil, nil, nil, false, indent, ""}
		if input[0] != '%' && input[0] != '#' && input[0] != '.' {
			err = &ParseError{1, 1}
			return
		}
		if input[len(input)-1] == '\n' {
			self.tag.LineBreak = "\n"
			input = input[0 : len(input)-1]
		}
		if len(indent) > 0 {
			self.tag.LineBreak = "\n"
		}
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
					self.tag.Name = string(input[start:i])
					break
				}
				if i == len(input)-1 {
					self.tag.Name = string(input[start : i+1])
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
					self.tag.Classes = append(self.tag.Classes, class)
					break
				}
				if i == len(input)-1 {
					class := string(input[start : i+1])
					self.tag.Classes = append(self.tag.Classes, class)
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
					self.tag.Id = string(input[start:i])
					break
				}
				if i == len(input)-1 {
					self.tag.Id = string(input[start : i+1])
				}
			}
			i -= 1
			continue
		}
		if input[i] == '(' || !self.completed {
			if !self.completed {
				i -= 1
			}
			var sn *StaticNode
			self.completed = false
			start = i + 1
			attr := Attribute{}
			for i = i + 1; i < len(input); i += 1 {
				switch input[i] {
				case '=':
					attr.Name = string(input[start:i])
					start = i + 1
				case ' ':
					sn = &StaticNode{}
					sn.Content = string(input[start+1 : i-1])
					attr.Value = sn
					self.tag.Attrs = append(self.tag.Attrs, attr)
					attr = Attribute{}
					start = i + 1
				case ')':
					self.completed = true
					sn = &StaticNode{}
					sn.Content = string(input[start+1 : i-1])
					attr.Value = sn
					goto EndAttrs
				}
			}
			sn = &StaticNode{}
			sn.Content = string(input[start+1 : i-1])
			attr.Value = sn
		EndAttrs:
			self.tag.Attrs = append(self.tag.Attrs, attr)
		}
		if i < len(input) && unicode.IsSpace(input[i]) {
			staticContent := string(input[i+1:])
			sn := &StaticNode{}
			sn.Content = staticContent
			self.tag.AddChild(sn)
			break
		}
		if i < len(input) && input[i] == '/' {
			self.tag.Close = true
		}
	}

	n = self.tag
	completed = self.completed
	return
}

type StaticParser struct {
}

func (self *StaticParser) Parse(indent string, input []rune) (n Node, completed bool, err *ParseError) {
	sn := &StaticLineNode{}
	sn.Content = string(input)
	sn.Indent = indent
	n = sn
	completed = true
	return
}

type DoctypeNode struct {
	Specifier string
}

func (self *DoctypeNode) Accept(visitor NodeVisitor) {
	visitor.VisitDoctype(self)
}

func (self *DoctypeNode) AddChild(child Node) (ok bool) {
	return
}

type TagNode struct {
	Name      string
	Id        string
	Classes   []string
	Attrs     []Attribute
	Children  []Node
	Close     bool
	Indent    string
	LineBreak string
}

func (self *TagNode) Accept(visitor NodeVisitor) {
	visitor.VisitTag(self)
}

func (self *TagNode) AddChild(child Node) (ok bool) {
	ok = child != nil
	self.Children = append(self.Children, child)
	return
}

type StaticNode struct {
	Content string
}

func (self *StaticNode) Accept(visitor NodeVisitor) {
}

func (self *StaticNode) AddChild(child Node) (ok bool) {
	return
}

type StaticLineNode struct {
	StaticNode
	Indent string
}

func (self *StaticLineNode) Accept(visitor NodeVisitor) {
	visitor.VisitStaticLine(self)
}

func (self *StaticLineNode) AddChild(child Node) (ok bool) {
	return
}
