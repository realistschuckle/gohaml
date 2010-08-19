package gohaml

import (
	"unicode"
	"strings"
)

type hamlParser struct {
}

func (self *hamlParser) parse(input string) (output *tree) {
	output = newTree()
	var currentNode *node
	j := 0
	for i, r := range input {
		if r == '\n' {
			node := parseLeadingSpace(input[j:i])
			if currentNode == nil  || node.indentLevel == currentNode.indentLevel {
				output.nodes.Push(node)
			} else if node.indentLevel > currentNode.indentLevel {
				currentNode.children.Push(node)
			}
			currentNode = node
			j = i + 1
		}
	}
	node := parseLeadingSpace(input[j:])
	if currentNode == nil  || node.indentLevel == currentNode.indentLevel {
		output.nodes.Push(node)
	} else if node.indentLevel > currentNode.indentLevel {
		currentNode.children.Push(node)
	}
	return
}

var parser hamlParser

func parseLeadingSpace(input string) (output *node) {
	node := new(node)
	for i, r := range input {
		switch {
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
			output.indentLevel = i
			break
		}
	}
	return
}

func parseKey(input string, node *node) (output *node) {
	if input[len(input) - 1] == '<' {
		output = parseNoNewline("", node)
		output.remainder.value = input[0:len(input) - 1]
	} else {
		output = node
		output.remainder.value = input
	}
	output.remainder.needsResolution = true
	return
}

func parseTag(input string, node *node) (output *node) {
	for i, r := range input {
		switch {
		case r == '.':
			output = parseClass(input[i + 1:], node)
		case r == '#':
			output = parseId(input[i + 1:], node)
		case r == '{':
			output = parseAttributes(tl(input[i + 1:]), node)
		case unicode.IsSpace(r):
			output = parseRemainder(input[i + 1:], node)
		}
		if nil != output {
			node.name = input[0:i]
			break
		}
	}
	if nil == output {
		node.name = input
		output = node;
	}
	return
}

func parseAttributes(input string, node *node) (output *node) {
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
			output = node
			break
		}
	}
	return
}

func parseId(input string, node *node) (output *node) {
	for i, r := range input {
		switch{
		case r == '.':
			node.addAttrNoLookup("id", input[0:i])
			output = parseClass(input[i + 1:], node)
		case r == '=':
			node.addAttrNoLookup("id", input[0:i])
			output = parseKey(tl(input[i + 1:]), node)
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

func parseClass(input string, node *node) (output *node) {
	for i, r := range input {
		switch {
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
		output = node
		output.addAttrNoLookup("class", input)
	}
	return
}

func parseRemainder(input string, node *node) (output *node) {
	if input[len(input) - 1] == '<' {
		output = parseNoNewline("", node)
		output.remainder.value = input[0:len(input) - 1]
	} else {
		output = node
		output.remainder.value = input
	}
	output.remainder.needsResolution = false
	return
}

func parseNoNewline(input string, node *node) (output *node) {
	node.noNewline = true
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
