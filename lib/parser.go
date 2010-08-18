package gohaml

import (
	"unicode"
)

type hamlParser struct {
}

func (self *hamlParser) parse(input string) (output *tree) {
	output = newTree()
	node := parseLeadingSpace(input)
	output.nodes.Push(node)
		
	return
}

func (self *hamlParser) initialize() {

}

var parser hamlParser

func init() {
}

func parseLeadingSpace(input string) (output *node) {
	node := new(node)
	for i, r := range input {
		switch {
		case r == '%':
			output = parseTag(input[i + 1:], node)
		case r == '.':
			output = parseClass(input[i + 1:], node)
		case !unicode.IsSpace(r):
			output = parseRemainder(input[i:], node)
		}
		if nil != output {break}
	}
	return
}

func parseTag(input string, node *node) (output *node) {
	for i, r := range input {
		switch {
		case r == '.':
			output = parseClass(input[i + 1:], node)
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

func parseClass(input string, node *node) (output *node) {
	for i, r := range input {
		switch {
		case r == '.':
			node.attrs.Push(&resPair{res{"class", false}, res{input[0:i], false}})
			output = parseClass(input[i + 1:], node)
		case unicode.IsSpace(r):
			node.attrs.Push(&resPair{res{"class", false}, res{input[0:i], false}})
			output = parseRemainder(input[i + 1:], node)
		}
		if nil != output {
			break;
		}
	}
	if nil == output {
		output = node
		node.attrs.Push(&resPair{res{"class", false}, res{input, false}})
	}
	return
}

func parseRemainder(input string, node *node) (output *node) {
	output = node
	output.remainder.value = input
	output.remainder.needsResolution = false
	return
}
