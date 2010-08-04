package gohaml

import (
	"bytes"
	"unicode"
	"fmt"
	"strings"
)

type attrMap map[string]string

func (m attrMap) String() (s string) {
	buf := bytes.NewBuffer([]byte{})
	for k, v := range m {
		buf.WriteString(" ")
		buf.WriteString(k)
		buf.WriteString("=\"")
		buf.WriteString(v)
		buf.WriteString("\"")
	}
	s = buf.String()
	return
}

type Engine struct {
	Options map[string]interface{}
	Indentation string
	input string
	parsingState int
}

const (
	leadingSpace = iota
	tagName
	keyName
	className
	content
	id
)

func NewEngine(input string) (engine *Engine) {
	engine = &Engine{make(map[string]interface{}), "\t", input, leadingSpace}
	engine.Options["autoclose"] = true
	return
}

func (engine *Engine) Render(scope map[string]interface{}) (output string) {
	autoclose := ""
	if engine.Options["autoclose"].(bool) {
		autoclose = " /"
	}
	lines := strings.Split(engine.input, "\n", -1)
	if 0 == len(lines) {lines = []string{engine.input}}
	for _, line := range lines {
		if 0 == len(line) {continue}
		_, tag, attrs, remainder := engine.splitTag(line, scope)
		switch len(tag) {
		case 0:
			output += remainder
		default:
			switch len(remainder) {
			case 0:
				output += fmt.Sprintf("<%s%s%s>", tag, attrs, autoclose)
			default:
				output += fmt.Sprintf("<%s%s>%s</%s>", tag, attrs, remainder, tag)
			}
		}
	}
	return
}

func (engine *Engine) splitTag(line string, scope map[string]interface{}) (indentCount int, tag string, attrs attrMap, remainder string) {
	i, r := 0, 0
	key := ""
	attrs = make(map[string]string)
	for i, r = range line {
		if '\\' == r && leadingSpace == engine.parsingState {
			remainder = line[1:]
			break
		}
		if '=' == r {
			engine.parsingState = keyName
			key = line[i + 1:]
			break
		}
		if '%' == r {
			engine.parsingState = tagName
			continue
		}
		if '#' == r {
			if 0 == len(tag) {tag = "div"}
			engine.parsingState = id
			continue
		}
		if '.' == r {
			if 0 == len(tag) {tag = "div"}
			if className == engine.parsingState {attrs["class"] += " "}
			engine.parsingState = className
			continue
		}
		if !unicode.IsSpace(r) && leadingSpace == engine.parsingState {
			remainder = line[i:]
			engine.parsingState = content
			break
		}
		if unicode.IsSpace(r) {
			if keyName == engine.parsingState {
				remainder = fmt.Sprint(scope[key])
			}
			engine.parsingState = content
			remainder += line[i:]
			break
		}
		switch engine.parsingState {
		case id:
			if _, ok := attrs["id"]; !ok {
				attrs["id"] = ""
			}
			attrs["id"] += line[i:i + 1]
		case tagName:
			tag += line[i:i + 1]
		case className:
			if _, ok := attrs["class"]; !ok {
				attrs["class"] = ""
			}
			attrs["class"] += line[i:i + 1]
		}
	}
	if keyName == engine.parsingState {
		key = strings.TrimLeftFunc(key, unicode.IsSpace)
		remainder = fmt.Sprint(scope[key])
	}
	tag = strings.TrimRightFunc(tag, unicode.IsSpace)
	remainder = strings.TrimLeftFunc(remainder, unicode.IsSpace)
	return
}
