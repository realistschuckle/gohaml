package gohaml

import (
	"unicode"
	"strings"
	"scanner"
	"fmt"
	"strconv"
	"os"
)

type hamlParser struct {
}

func (self *hamlParser) parse(input string) (output *tree) {
	output = newTree()
	var currentNode inode
	j := 0
	for i, r := range input {
		if r == '\n' {
			node := parseLeadingSpace(input[j:i])
			putNodeInPlace(currentNode, node, output)
			currentNode = node
			j = i + 1
		}
	}
	node := parseLeadingSpace(input[j:])
	putNodeInPlace(currentNode, node, output)
	return
}

func putNodeInPlace(cn inode, node inode, t *tree) {
	if cn == nil || cn.nil() {
		t.nodes.Push(node)
	} else if node.indentLevel() < cn.indentLevel() {
		for cn = cn.parent(); cn != nil && node.indentLevel() < cn.indentLevel(); cn = cn.parent() {}
		putNodeInPlace(cn, node, t)
	} else if node.indentLevel() == cn.indentLevel() && cn.parent() != nil{
		cn.parent().addChild(node)
	} else if node.indentLevel() == cn.indentLevel() {
		t.nodes.Push(node)
	} else if node.indentLevel() > cn.indentLevel() {
		cn.addChild(node)
	}
}

var parser hamlParser

func parseLeadingSpace(input string) (output inode) {
	node := new(node)
	for i, r := range input {
		switch {
		case r == '-':
			output = parseCode(input[i + 1:], node)
		case r == '%':
			output = parseTag(input[i + 1:], node)
		case r == '#':
			output = parseId(input[i + 1:], node)
		case r == '.':
			output = parseClass(input[i + 1:], node)
		case r == '=':
			output = parseKey(tl(input[i + 1:]), node)
		case r == '\\':
			output = parseRemainder(input[i + 1:], node)
		case !unicode.IsSpace(r):
			output = parseRemainder(input[i:], node)
		}
		if nil != output {
			output.setIndentLevel(i)
			break
		}
	}
	return
}

func parseKey(input string, n *node) (output inode) {
	if input[len(input) - 1] == '<' {
		n = parseNoNewline("", n)
		n.setRemainder(input[0:len(input) - 1], true)
	} else {
		n.setRemainder(input, true)
		output = n
	}
	output = n
	return
}

func parseTag(input string, node *node) (output inode) {
	for i, r := range input {
		switch {
		case r == '.':
			output = parseClass(input[i + 1:], node)
		case r == '#':
			output = parseId(input[i + 1:], node)
		case r == '{':
			output = parseAttributes(tl(input[i + 1:]), node)
		case r == '<':
			output = parseNoNewline(input[i + 1:], node)
		case r == '=':
			output = parseKey(tl(input[i + 1:]), node)
		case r == '/':
			output = parseAutoclose("", node)
		case unicode.IsSpace(r):
			output = parseRemainder(input[i + 1:], node)
		}
		if nil != output {
			node._name = input[0:i]
			break
		}
	}
	if nil == output {
		node._name = input
		output = node;
	}
	return
}

func parseAutoclose(input string, node *node) (output inode) {
	node._autoclose = true
	output = node
	return
}

func parseAttributes(input string, node *node) (output inode) {
	inKey := true
	inRocket := false
	keyEnd, attrStart := 0, 0
	for i, r := range input {
		if inKey && (r == '=' || unicode.IsSpace(r)) {
			inKey = false
			inRocket = true
			keyEnd = i
		} else if inRocket && r != '>' && r != '=' && !unicode.IsSpace(r) {
			inRocket = false
			attrStart = i
		} else if r == ',' {
			node.addAttr(t(input[0:keyEnd]), t(input[attrStart:i]))
			output = parseAttributes(tl(input[i + 1:]), node)
			break
		} else if r == '}' {
			node.addAttr(input[0:keyEnd], t(input[attrStart:i]))
			output = parseTag(input[i + 1:], node)
			break
		}
	}
	return
}

func parseId(input string, node *node) (output inode) {
	for i, r := range input {
		switch{
		case r == '.':
			node.addAttrNoLookup("id", input[0:i])
			output = parseClass(input[i + 1:], node)
		case r == '=':
			node.addAttrNoLookup("id", input[0:i])
			output = parseKey(tl(input[i + 1:]), node)
		case r == '{':
			node.addAttrNoLookup("id", input[0:i])
			output = parseAttributes(tl(input[i + 1:]), node)
		case unicode.IsSpace(r):
			node.addAttrNoLookup("id", input[0:i])
			output = parseRemainder(input[i + 1:], node)
		}
		if nil != output {break}
	}
	if nil == output {
		output = node
		node.addAttrNoLookup("id", input)
	}
	return
}

func parseClass(input string, node *node) (output inode) {
	for i, r := range input {
		switch {
		case r == '{':
			node.addAttrNoLookup("class", input[0:i])
			output = parseAttributes(tl(input[i + 1:]), node)
		case r == '.':
			node.addAttrNoLookup("class", input[0:i])
			output = parseClass(input[i + 1:], node)
		case r == '=':
			node.addAttrNoLookup("class", input[0:i])
			output = parseKey(tl(input[i + 1:]), node)
		case unicode.IsSpace(r):
			node.addAttrNoLookup("class", input[0:i])
			output = parseRemainder(input[i + 1:], node)
		}
		if nil != output {break}
	}
	if nil == output {
		node.addAttrNoLookup("class", input)
		output = node
	}
	return
}

func parseRemainder(input string, node *node) (output inode) {
	if input[len(input) - 1] == '<' {
		node = parseNoNewline("", node)
		node._remainder.value = input[0:len(input) - 1]
		node._remainder.needsResolution = false
	} else {
		node._remainder.value = input
		node._remainder.needsResolution = false
	}
	output = node
	return
}

func parseNoNewline(input string, node *node) (output *node) {
	node.setNoNewline(true)
	output = node
	return
}

func t(input string) (output string) {
	output = strings.Trim(input, " 	")
	return
}

func tl(input string) (output string) {
 	output = strings.TrimLeft(input, " 	")
	return
}

const eof int = yyUSER + 1

func parseCode(input string, node inode) (output inode) {
	s.Init(strings.NewReader(input))
	success, result := yyparse(eof, scan)
	if !success {
		fmt.Fprintf(os.Stderr, "Did not recognize %s", input)
	}
	output = result
	return
}

var s scanner.Scanner

func scan(v *yystype) (output int) {
	i := s.Scan()
	switch i {
	case scanner.Ident:
		output = ident
		v.s = s.TokenText()
	case scanner.String, scanner.RawString:
		output = atom
		text := s.TokenText()
		v.i = text[1:len(text) - 1]
	case scanner.Int:
		output = atom
		v.i, _ = strconv.Atoi(s.TokenText())
	case scanner.Float:
		output = atom
		v.i, _ = strconv.Atof(s.TokenText())
	case scanner.EOF:
		output = eof
	default:
		output = i
	}
	return
}

