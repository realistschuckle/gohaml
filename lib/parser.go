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
	indentation string   // one indentation level
	filter      int      // the indent level at which to filter.
	filterNode  *node    // the node filtered content belongs to.
	filterbuff  []string // a buffer for lines of filtered content.
	FilterMap            // a dictionary of filters.
}

func newHamlParser(indentation string) *hamlParser {
	return &hamlParser{indentation, -1, nil, nil, defaultFilterMap}
}

// Append a line of filter content to self.filterbuff with a specified
// indentation depth (number of tabs).
func (self *hamlParser) appendFiltered(indent int, content string) {
	// If an inline value was supplied, indent it and append to the filterbuff.
	self.filterbuff = append(self.filterbuff,
		fmt.Sprintf("%s%s",
			strings.Repeat(self.indentation, indent),
			strings.TrimLeftFunc(content, unicode.IsSpace)))
}
// Pass the contents of filterNode (not the struct field) to the filter named
// in the node. Any error due to a missing/badly-named filter will be returned.
func (self *hamlParser) wrapFilter(filterNode *node, line int) os.Error {
	// Search for the named filter.
	name := filterNode._name[1:]
	fn, found := self.FilterMap[name]
	filterNode._name = ""
	if !found {
		return fmt.Errorf("Line %d: Filter not found %s", line, name)
	}

	// Compute filter input
	input := strings.Join(self.filterbuff, "")
	self.filterbuff = nil
	if input == "" || input[len(input)-1] != '\n' {
		input += "\n"
	}

	// Compute filter indentation amount and reset self.filter
	var indent string
	if self.filter > 0 {
		indent = strings.Repeat(self.indentation, self.filter)
	}
	self.filter = -1

	// Compute filtered output.
	filterNode._remainder.value = fn.Filter(input[:len(input)-1], indent, self.indentation)
	return nil
}
// Check if filtering was/is being performed. Compute filtered output for
// completed filted blocks. Append new filtered content to the filter buffer.
// Keep the content of the filtered blocks out of the tree.
func (self *hamlParser) filterInput(n inode, input string, filtering bool, line int) (in inode, err os.Error) {
	in = n
	var nod *node
	switch n.(type) {
	case nil:
	case *node:
		nod = n.(*node)
	default:
		return
	}
	if nod != nil && self.filter >= 0 && len(nod._name) > 0 { // A filter terminated with a new filter.
		if err = self.wrapFilter(self.filterNode, line); err != nil {
			return
		}
	}
	switch {
	case filtering && self.filter < 0:
		// The parser just hit a :filter (the node should not be nil)
		self.filterNode, self.filter = nod, nod.indentLevel()
		if n := nod; len(n._remainder.value) > 0 {
			self.appendFiltered(n.indentLevel()+1, n._remainder.value)
		}
	case filtering:
		// Indent the last line and append it to the filter (with a newline byte)
		switch indent := self.filter; {
		case nod != nil:
			indent = nod.indentLevel()
			in = nil
			fallthrough
		default:
			self.appendFiltered(indent, input)
		}
	case self.filter >= 0: // We were filtering, but now out of filter scope.
		var old *node
		old, self.filterNode = self.filterNode, nil
		if err = self.wrapFilter(old, line); err != nil {
			return
		}
	}
	return
}
// Preform the same function as filterInput, but at the end of the input string.
// Behaves differently from filterInput because the filtered output is always
// computed at the end.
func (self *hamlParser) filterLast(n inode, input string, filtering bool, line int) (in inode, err os.Error) {
	in = n
	var nod *node
	switch n.(type) {
	case nil:
	case *node:
		nod = n.(*node)
	default:
		return
	}
	if nod != nil && self.filter >= 0 && len(nod._name) > 0 { // A filter terminated with a new filter.
		if err = self.wrapFilter(self.filterNode, line); err != nil {
			return
		}
	}
	// Prepare filterNode and call wrapFilter.
	switch {
	case filtering && self.filter >= 0: // Parse was filtering before the last line, and the last line was a filter.
		if nod != nil {
			self.appendFiltered(nod.indentLevel(), input)
		}
		fallthrough
	case filtering:
		// Check self.filter in case of fallthough.
		if self.filter < 0 { // The last line contains a :filter (and the node is not nil). 
			self.filterNode = nod
			if n := nod; len(n._remainder.value) > 0 {
				self.appendFiltered(n.indentLevel()+1, n._remainder.value)
			}
		}
		fallthrough
	case self.filter >= 0:
		// In all cases, wrap the new filter because there is no more input.
		var old *node
		old, self.filterNode = self.filterNode, nil
		if err = self.wrapFilter(old, line); err != nil {
			return
		}
	}
	return
}

func (self *hamlParser) parse(input string) (output *tree, err os.Error) {
	output = newTree()
	var currentNode inode
	var nod inode
	var filtering bool
	lastSpaceChar := -1
	line := 1
	j := 0
	for i, r := range input {
		if r == '\n' {
			line += 1
			nod, err, lastSpaceChar, filtering = parseLeadingSpace(input[j:i], lastSpaceChar, line, self.filter)
			if err != nil {
				return
			}
			// Filter the input if necessary, and possibly skip added nod to the tree.
			nod, err = self.filterInput(nod, input[j:i+1], filtering, line)
			if err != nil {
				return
			}
			if nod != nil && !nod.nil() {
				putNodeInPlace(currentNode, nod, output)
				currentNode = nod
			}
			j = i + 1
		}
	}
	nod, err, lastSpaceChar, filtering = parseLeadingSpace(input[j:], lastSpaceChar, line, self.filter)
	if err != nil {
		return
	}
	// Filter the input if necessary.
	nod, err = self.filterLast(nod, input[j:], filtering, line)
	if err != nil {
		return
	}
	if nod != nil && !nod.nil() {
		putNodeInPlace(currentNode, nod, output)
	}
	return
}

func putNodeInPlace(cn inode, node inode, t *tree) {
	if node == nil || node.nil() {
		return
	}
	if cn == nil || cn.nil() {
		t.nodes.Push(node)
	} else if node.indentLevel() < cn.indentLevel() {
		for cn = cn.parent(); cn != nil && node.indentLevel() < cn.indentLevel(); cn = cn.parent() {
		}
		putNodeInPlace(cn, node, t)
	} else if node.indentLevel() == cn.indentLevel() && cn.parent() != nil {
		cn.parent().addChild(node)
	} else if node.indentLevel() == cn.indentLevel() {
		t.nodes.Push(node)
	} else if node.indentLevel() > cn.indentLevel() {
		cn.addChild(node)
	}
}

var parser = newHamlParser("\t")

func parseLeadingSpace(input string, lastSpaceChar int, line int, filter int) (output inode, err os.Error, spaceChar int, inFilter bool) {
	nod := new(node)
	for i, r := range input {
		switch {
		case filter >= 0 && i > filter && !unicode.IsSpace(r):
			// Now inside the filter's scope. 
			output = parseRemainderCDATA(input[i:], nod, line)
			inFilter = true
		case r == ':':
			// Found a filter instantiation.
			output, err = parseFilter(input[i:], nod, line) // filter name has ':'
			inFilter = true
		case r == '-':
			output = parseCode(input[i+1:], nod, line)
		case r == '%':
			output, err = parseTag(input[i+1:], nod, true, line)
		case r == '#':
			output, err = parseId(input[i+1:], nod, line)
		case r == '.':
			output, err = parseClass(input[i+1:], nod, line)
		case r == '=':
			output = parseKey(tl(input[i+1:]), nod, line)
		case r == '\\':
			output = parseRemainder(input[i+1:], nod, line)
		case !unicode.IsSpace(r):
			output = parseRemainder(input[i:], nod, line)
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
	if output == nil {
		// Keep the inFilter value accurate on blank lines.
		inFilter = filter >= 0
	}
	spaceChar = lastSpaceChar
	return
}

// Parse a filter instantiation ":filter [content]".
// The filter is stored as a regular *node. The _name is the name of the filter.
// The _remainder.value is any content appearing on the same line as the filter.
// Filtered content is always kept as plain text, and never parsed as haml.
func parseFilter(input string, n *node, line int) (output inode, err os.Error) {
	if len(input) == 0 {
		err = os.NewError(fmt.Sprintf("Parse error on line %d: Empty input\n", line))
		return
	}
	if len(input) == 1 {
		err = os.NewError(fmt.Sprintf("Syntax error on line %d: Missing filter name\n", line))
		return
	}
	for i, r := range input {
		if unicode.IsSpace(r) {
			n._name = input[:i]
			if len(n._name) <= 1 {
				err = os.NewError(fmt.Sprintf("Syntax error on line %d: Missing filter name\n", line))
				return
			}
			n.setRemainder(input[i:]+"\n", false)
			output = n
			return
		}
	}
	n._name = input
	output = n
	return
}

func parseKey(input string, n *node, line int) (output inode) {
	if input[len(input)-1] == '<' {
		n = parseNoNewline("", n, line)
		n.setRemainder(input[0:len(input)-1], true)
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
			output, err = parseClass(input[i+1:], node, line)
		case r == '#':
			output, err = parseId(input[i+1:], node, line)
		case r == '{':
			output, err = parseAttributes(tl(input[i+1:]), node, line)
		case r == '<':
			output = parseNoNewline(input[i+1:], node, line)
		case r == '=':
			output = parseKey(tl(input[i+1:]), node, line)
		case r == '/':
			output = parseAutoclose("", node, line)
		case unicode.IsSpace(r):
			output = parseRemainder(input[i+1:], node, line)
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
		output = node
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
			output, err = parseAttributes(tl(input[i+1:]), node, line)
			break
		} else if r == '}' {
			if attrStart == 0 {
				msg := fmt.Sprintf("Syntax error on line %d: Attribute requires a value.\n", line)
				err = os.NewError(msg)
				return
			}
			if inKey {
				msg := fmt.Sprintf("Syntax error on line %d: Attribute requires a rocket and value.\n", line)
				err = os.NewError(msg)
				return
			}
			attrValue := t(input[attrStart:i])
			node.addAttr(input[0:keyEnd], attrValue)
			output, _ = parseTag(input[i+1:], node, false, line)
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
		switch {
		case r == '.':
			output, _ = parseClass(input[i+1:], node, line)
		case r == '=':
			output = parseKey(tl(input[i+1:]), node, line)
		case r == '{':
			output, err = parseAttributes(tl(input[i+1:]), node, line)
		case unicode.IsSpace(r):
			output = parseRemainder(input[i+1:], node, line)
		}
		if nil != output {
			break
		}
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
			output, err = parseAttributes(tl(input[i+1:]), node, line)
		case r == '.':
			output, err = parseClass(input[i+1:], node, line)
		case r == '=':
			output = parseKey(tl(input[i+1:]), node, line)
		case unicode.IsSpace(r):
			output = parseRemainder(input[i+1:], node, line)
		}
		if nil != output {
			break
		}
	}
	if nil == output {
		node.addAttrNoLookup("class", input)
		output = node
	}
	return
}

// Just chunk up the rest of the line. -bmatsuo
func parseRemainderCDATA(input string, node *node, line int) (output inode) {
	node._remainder.value = input
	node._remainder.needsResolution = false
	output = node
	return
}

func parseRemainder(input string, node *node, line int) (output inode) {
	if input[len(input)-1] == '<' {
		node = parseNoNewline("", node, line)
		node._remainder.value = input[0 : len(input)-1]
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
		v.i = text[1 : len(text)-1]
	case scanner.Int:
		output = atom
		v.i, _ = strconv.Atoi(s.TokenText())
	case scanner.Float:
		output = atom
		v.i, _ = strconv.Atof64(s.TokenText())
	case scanner.EOF:
		output = eof
	default:
		output = i
	}
	return
}
