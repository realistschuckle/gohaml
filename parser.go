package gohaml

import (
	"bytes"
	"container/list"
	"io"
	"errors"
	"strings"
)

type parsedDoc struct {
	nodes *list.List
}

func (self *parsedDoc) Accept(c HamlCompiler) {
	for e := self.nodes.Front(); e != nil; e = e.Next() {
		n := e.Value.(Node)
		n.Accept(c)
	}
	return
}

type parser struct {
	input  *strings.Reader
	states []NodeParser
}

func (self *parser) Parse() (d ParsedDocument, e error) {
	if self.states == nil {
		e = errors.New("parser requires start states to Parse")
		return
	}
	states := self.states
	l := list.New()
	for self.input.Len() > 0 {
		for _, nodeParser := range states {
			var node Node
			if node, e = nodeParser.Parse(self.input); e == nil {
				l.PushFront(node)
				states = nodeParser.Next()
				goto cont
			}
		}
		e = errors.New("Did not fully parse the input")
		return

	cont:
	}
	d = &parsedDoc{l}
	return
}

type docTypeParser struct{}

func (self *docTypeParser) Parse(i *strings.Reader) (n Node, e error) {
	failed := false
	x := 0
	var offset int64 = 0
	for x = 0; !failed && x < 3; x += 1 {
		r, w, err := i.ReadRune()
		offset -= int64(w)
		failed = r != '!' || err != nil
	}
	if failed {
		i.Seek(offset, 1)
		e = errors.New("Could not parse input")
	} else {
		buf := bytes.Buffer{}
		for r, _, err := i.ReadRune(); err == nil; r, _, err = i.ReadRune() {
			if r == '\n' {
				i.UnreadRune()
				break
			}
			buf.WriteRune(r)
		}
		n = &DocTypeNode{strings.TrimLeft(buf.String(), " ")}
	}
	return
}

func (self *docTypeParser) Next() (n []NodeParser) {
	return
}

type tagParser struct{}

func (self *tagParser) Parse(i *strings.Reader) (n Node, e error) {
	var r rune
	if r, _, e = i.ReadRune(); e != nil {
		return
	}

	if r != '%' {
		i.UnreadRune()
		e = errors.New("Not a tag")
		return
	}

	b := bytes.Buffer{}
	forceClose := false
	cont := true
	for r, _, e = i.ReadRune(); e == nil && cont; r, _, e = i.ReadRune() {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '_' || r == ':':
			b.WriteRune(r)
		case r == '/':
			forceClose = true
			cont = false
		default:
			i.UnreadRune()
			cont = false
		}
	}

	if e == nil {
		i.UnreadRune()
	}
	if e == io.EOF || forceClose {
		e = nil
	}

	n = &TagNode{b.String(), forceClose}

	return
}

func (self *tagParser) Next() (n []NodeParser) {
	n = []NodeParser{&classNameParser{}}
	return
}

type classNameParser struct{}

func (self *classNameParser) Parse(i *strings.Reader) (n Node, e error) {
	var r rune
	if r, _, e = i.ReadRune(); e != nil {
		return
	}

	if r != '.' {
		i.UnreadRune()
		e = errors.New("Not a class name")
		return
	}

	b := bytes.Buffer{}
	cont := true
	for r, _, e = i.ReadRune(); e == nil && cont; r, _, e = i.ReadRune() {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '_':
			b.WriteRune(r)
		default:
			cont = false
		}
	}

	if e == io.EOF {
		e = nil
	} else {
		i.UnreadRune()
	}

	n = &ClassNameNode{b.String()}

	return
}

func (self *classNameParser) Next() (n []NodeParser) {
	return
}