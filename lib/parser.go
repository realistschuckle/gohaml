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

func (self *hamlParser) parse(input string) (output *tree, err os.Error) {
	output = newTree()
	var currentNode inode
	var node inode
	lastSpaceChar := -1
	line := 1
	j := 0
	for i, r := range input {
		if r == '\n' {
			line += 1
			node, err, lastSpaceChar = parseLeadingSpace(input[j:i], lastSpaceChar, line)
			if err != nil {
				return
			}
			if node != nil && !node.nil() {
				putNodeInPlace(currentNode, node, output)
				currentNode = node
			}
			j = i + 1
		}
	}
	node, err, lastSpaceChar = parseLeadingSpace(input[j:], lastSpaceChar, line)
	if err != nil {
		return
	}
	if node != nil && !node.nil() {
		putNodeInPlace(currentNode, node, output)
	}
	return
}

func putNodeInPlace(cn inode, node inode, t *tree) {
	if node == nil || node.nil() {return}
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

func parseLeadingSpace(input string, lastSpaceChar int, line int) (output inode, err os.Error, spaceChar int) {
	node := new(node)
	for i, r := range input {
		switch {
		case r == '-':
			output = parseCode(input[i + 1:], node, line)
		case r == '%':
			output, err = parseTag(input[i + 1:], node, true, line)
		case r == '#':
			output, err = parseId(input[i + 1:], node, line)
		case r == '.':
			output, err = parseClass(input[i + 1:], node, line)
		case r == '=':
			output = parseKey(tl(input[i + 1:]), node, line)
		case r == '\\':
			output = parseRemainder(input[i + 1:], node, line)
		case !unicode.IsSpace(r):
			output = parseRemainder(input[i:], node, line)
		case unicode.IsSpace(r):
			if lastSpaceChar > 0 && r != lastSpaceChar {
				from, to := "space", "tab"
				if lastSpaceChar == 9 {
					from = "tab"
					to = "space"
				}
				msg := fmt.Sprintf("Syntax error on line %d: Inconsistent spacing in document changed from %s to %s characters.\n", line, from, to)
				err = os.NewError(msg)
			} else {
				lastSpaceChar = r
			}
		}
		if nil != err {
			break
		}
		if nil != output {
			output.setIndentLevel(i)
			break
		}
	}
	spaceChar = lastSpaceChar
	return
}

func parseKey(input string, n *node, line int) (output inode) {
	if input[len(input) - 1] == '<' {
		n = parseNoNewline("", n, line)
		n.setRemainder(input[0:len(input) - 1], true)
	} else {
		n.setRemainder(input, true)
		output = n
	}
	output = n
	return
}

func parseTag(input string, node *node, newTag bool, line int) (output inode, err os.Error) {
	if 0 == len(input) && newTag {
		err = os.NewError(fmt.Sprintf("Syntax error on line %d: Invalid tag: %s.\n", line, input))
		return
	}
	for i, r := range input {
		switch {
		case r == '.':
			output, err = parseClass(input[i + 1:], node, line)
		case r == '#':
			output, err = parseId(input[i + 1:], node, line)
		case r == '{':
			output, err = parseAttributes(tl(input[i + 1:]), node, line)
		case r == '<':
			output = parseNoNewline(input[i + 1:], node, line)
		case r == '=':
			output = parseKey(tl(input[i + 1:]), node, line)
		case r == '/':
			output = parseAutoclose("", node, line)
		case unicode.IsSpace(r):
			output = parseRemainder(input[i + 1:], node, line)
		}
		if nil != err {
			break
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

func parseAutoclose(input string, node *node, line int) (output inode) {
	node._autoclose = true
	output = node
	return
}

func parseAttributes(input string, node *node, line int) (output inode, err os.Error) {
	inKey := true
	inRocket := false
	keyEnd, attrStart := 0, 0
	for i, r := range input {
		if inKey && (r == '=' || unicode.IsSpace(r)) {
			inKey = false
			inRocket = true
			keyEnd = i
		} else if inRocket && r != '>' && r != '=' && r != '}' && !unicode.IsSpace(r) {
			inRocket = false
			attrStart = i
		} else if r == ',' {
			node.addAttr(t(input[0:keyEnd]), t(input[attrStart:i]))
			output, err = parseAttributes(tl(input[i + 1:]), node, line)
			break
		} else if r == '}' {
			if attrStart == 0 {
				msg := fmt.Sprintf("Syntax error on line %d: Attribute requires a value.\n", line)
				err = os.NewError(msg)
				return
			}
			if(inKey) {
				msg := fmt.Sprintf("Syntax error on line %d: Attribute requires a rocket and value.\n", line)
				err = os.NewError(msg)
				return
			}
			attrValue := t(input[attrStart:i])
			node.addAttr(input[0:keyEnd], attrValue)
			output, _ = parseTag(input[i + 1:], node, false, line)
			break
		}
	}
	if nil == output {
		msg := fmt.Sprintf("Syntax error on line %d: Attributes must have closing '}'.\n", line)
		err = os.NewError(msg)
	}
	return
}

func parseId(input string, node *node, line int) (output inode, err os.Error) {
	defer func() {
		if nil == output {
			msg := fmt.Sprintf("Syntax error on line %d: Illegal element: classes and ids must have values.\n", line)
			err = os.NewError(msg)
		}
	}()
	if len(input) == 0 {
		return
	}
	for i, r := range input {
		if r == '.' || r == '=' || r == '{' || unicode.IsSpace(r) {
			if i == 0 {
				return
			}
			node.addAttrNoLookup("id", input[0:i])			
		}
		switch{
		case r == '.':
			output, _ = parseClass(input[i + 1:], node, line)
		case r == '=':
			output = parseKey(tl(input[i + 1:]), node, line)
		case r == '{':
			output, err = parseAttributes(tl(input[i + 1:]), node, line)
		case unicode.IsSpace(r):
			output = parseRemainder(input[i + 1:], node, line)
		}
		if nil != output {break}
	}
	if nil == output {
		output = node
		node.addAttrNoLookup("id", input)
	}
	return
}

func parseClass(input string, node *node, line int) (output inode, err os.Error) {
	defer func() {
		if nil == output {
			msg := fmt.Sprintf("Syntax error on line %d: Illegal element: classes and ids must have values.\n", line)
			err = os.NewError(msg)
		}
	}()
	if len(input) == 0 {
		return
	}
	for i, r := range input {
		if r == '{' || r == '.' || r == '=' || unicode.IsSpace(r) {
			if i == 0 {
				return
			}
			node.addAttrNoLookup("class", input[0:i])
		}
		switch {
		case r == '{':
			output, err = parseAttributes(tl(input[i + 1:]), node, line)
		case r == '.':
			output, err = parseClass(input[i + 1:], node, line)
		case r == '=':
			output = parseKey(tl(input[i + 1:]), node, line)
		case unicode.IsSpace(r):
			output = parseRemainder(input[i + 1:], node, line)
		}
		if nil != output {break}
	}
	if nil == output {
		node.addAttrNoLookup("class", input)
		output = node
	}
	return
}

func parseRemainder(input string, node *node, line int) (output inode) {
	if input[len(input) - 1] == '<' {
		node = parseNoNewline("", node, line)
		node._remainder.value = input[0:len(input) - 1]
		node._remainder.needsResolution = false
	} else {
		node._remainder.value = input
		node._remainder.needsResolution = false
	}
	output = node
	return
}

func parseNoNewline(input string, node *node, line int) (output *node) {
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

func parseCode(input string, node inode, line int) (output inode) {
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
		switch s.TokenText() {
		case "for":
			output = tfor
		case "range":
			output = trange
		default:
			output = ident
		}
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

