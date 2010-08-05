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
	parsingState *state
	tag string
	attrs attrMap
	remainder string
	indentCount int
}

func NewEngine(input string) (engine *Engine) {
	engine = &Engine{make(map[string]interface{}), "\t", input, nil, "", make(map[string]string), "", 0}
	engine.makeStates()
	engine.Options["autoclose"] = true
	return
}

func (self *Engine) makeStates() {
	exitLeadingSpace := func(s *state, line []int, scope map[string]interface{}) {
		self.indentCount++
		return
	}
	
	exitTagState := func(s *state, line []int, scope map[string]interface{}) {
		self.tag = string(line[s.leftIndex + 1:s.rightIndex])
		return
	}
	
	exitIdState := func(s *state, line []int, scope map[string]interface{}) {
		if 0 == len(self.tag) {self.tag = "div"}
		self.attrs["id"] = string(line[s.leftIndex + 1:s.rightIndex])
		return
	}
	
	exitClassState := func(s *state, line []int, scope map[string]interface{}) {
		if 0 == len(self.tag) {self.tag = "div"}
		if _, ok := self.attrs["class"]; !ok {
			self.attrs["class"] = string(line[s.leftIndex + 1:s.rightIndex])
		} else {
			self.attrs["class"] += " " + string(line[s.leftIndex + 1:s.rightIndex])
		}
	}
	
	exitKeyState := func(s *state, line []int, scope map[string]interface{}) {
		key := string(line[s.leftIndex + 1:s.rightIndex])
		key = strings.TrimFunc(key, unicode.IsSpace)
		self.remainder = fmt.Sprint(scope[key])
	}
	
	exitContentState := func(s *state, line []int, scope map[string]interface{}) {
		self.remainder = string(line[s.leftIndex:s.rightIndex])
		return
	}
	
	exitBackslashState := func(s *state, line []int, scope map[string]interface{}) {
		self.remainder = string(line[1:])
		return
	}
	
 	leadingSpaceState := newState(exitLeadingSpace)
	tagState := newState(exitTagState)
	idState := newState(exitIdState)
	classState := newState(exitClassState)
	keyState := newState(exitKeyState)
	contentState := newState(exitContentState)
	backslashState := newState(exitBackslashState)
	
	matchTag := func(rune int) bool {return '%' == rune}
	matchId := func(rune int) bool {return '#' == rune}
	matchClass := func(rune int) bool {return '.' == rune}
	matchKey := func(rune int) bool {return '=' == rune}
	matchLeadingContent := func(rune int) bool {return !unicode.IsSpace(rune)}
	matchContent := func(rune int) bool {return unicode.IsSpace(rune)}
	matchBackslashState := func(rune int) bool {return '\\' == rune && 0 == self.indentCount}
	
	leadingSpaceState.addTransition(matchBackslashState, backslashState)
	leadingSpaceState.addTransition(matchTag, tagState)
	leadingSpaceState.addTransition(matchId, idState)
	leadingSpaceState.addTransition(matchClass, classState)
	leadingSpaceState.addTransition(matchKey, keyState)
	leadingSpaceState.addTransition(matchLeadingContent, contentState)
	
	tagState.addTransition(matchContent, contentState)
	tagState.addTransition(matchClass, classState)
	tagState.addTransition(matchId, idState)
	
	idState.addTransition(matchClass, classState)
	idState.addTransition(matchKey, keyState)
	idState.addTransition(matchContent, contentState)
	
	classState.addTransition(matchClass, classState)
	classState.addTransition(matchKey, keyState)
	classState.addTransition(matchContent, contentState)
	
	self.parsingState = leadingSpaceState
	
	return
}

func (self *Engine) Render(scope map[string]interface{}) (output string) {
	autoclose := ""
	if self.Options["autoclose"].(bool) {
		autoclose = " /"
	}
	lines := strings.Split(self.input, "\n", -1)
	if 0 == len(lines) {lines = []string{self.input}}
	for _, line := range lines {
		if 0 == len(line) {continue}
		self.parseLine(line, scope)
		switch len(self.tag) {
		case 0:
			output += self.remainder
		default:
			switch len(self.remainder) {
			case 0:
				output += fmt.Sprintf("<%s%s%s>", self.tag, self.attrs, autoclose)
			default:
				output += fmt.Sprintf("<%s%s>%s</%s>", self.tag, self.attrs, self.remainder, self.tag)
			}
		}
	}
	return
}

func (self *Engine) String() string {
	return fmt.Sprintf("<Engine tag: %s\nattrs: %s\nremainder: %s>", self.tag, self.attrs, self.remainder)
}

func (self *Engine) parseLine(line string, scope map[string]interface{}) {
	linen := []int(line)
	for i, _ := range linen {
		self.parsingState = self.parsingState.input(i, linen, scope)
	}
	self.parsingState.exit(self.parsingState, linen, scope)
	self.tag = strings.TrimRightFunc(self.tag, unicode.IsSpace)
	self.remainder = strings.TrimLeftFunc(self.remainder, unicode.IsSpace)
	return
}
