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
	startState *state
	tag string
	attrs attrMap
	remainder string
	indentCount int
	closeTag bool
}

func NewEngine(input string) (engine *Engine) {
	engine = &Engine{make(map[string]interface{}), "\t", input, nil, nil, "", make(map[string]string), "", 0, false}
	engine.makeStates()
	engine.Options["autoclose"] = true
	return
}

func (self *Engine) makeStates() {
	exitLeadingSpace := func(s *state, line []int, scope map[string]interface{}) {
		self.indentCount = s.rightIndex
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
	
	exitCloseTagState := func(s *state, line []int, scope map[string]interface{}) {
		self.closeTag = true
		return
	}
	
 	leadingSpaceState := newState(exitLeadingSpace)
	tagState := newState(exitTagState)
	idState := newState(exitIdState)
	classState := newState(exitClassState)
	keyState := newState(exitKeyState)
	contentState := newState(exitContentState)
	backslashState := newState(exitBackslashState)
	closeTagState := newState(exitCloseTagState)
	
	matchTag := func(rune int) bool {return '%' == rune}
	matchId := func(rune int) bool {return '#' == rune}
	matchClass := func(rune int) bool {return '.' == rune}
	matchKey := func(rune int) bool {return '=' == rune}
	matchLeadingContent := func(rune int) bool {return !unicode.IsSpace(rune)}
	matchContent := func(rune int) bool {return unicode.IsSpace(rune)}
	matchBackslashState := func(rune int) bool {return '\\' == rune && 0 == self.indentCount}
	matchCloseTagState := func(rune int) bool {return '/' == rune}
	
	leadingSpaceState.addTransition(matchBackslashState, backslashState)
	leadingSpaceState.addTransition(matchTag, tagState)
	leadingSpaceState.addTransition(matchId, idState)
	leadingSpaceState.addTransition(matchClass, classState)
	leadingSpaceState.addTransition(matchKey, keyState)
	leadingSpaceState.addTransition(matchLeadingContent, contentState)

	tagState.addTransition(matchCloseTagState, closeTagState)
	tagState.addTransition(matchClass, classState)
	tagState.addTransition(matchId, idState)
	tagState.addTransition(matchContent, contentState)
	
	idState.addTransition(matchClass, classState)
	idState.addTransition(matchKey, keyState)
	idState.addTransition(matchContent, contentState)
	
	classState.addTransition(matchClass, classState)
	classState.addTransition(matchKey, keyState)
	classState.addTransition(matchContent, contentState)
	
	self.parsingState = leadingSpaceState
	self.startState = leadingSpaceState
	
	return
}

func (self *Engine) Render(scope map[string]interface{}) (output string) {
	lines := strings.Split(self.input, "\n", -1)
	if 0 == len(lines) {lines = []string{self.input}}
	lineEnd := "\n"
	for i, line := range lines {
		if 0 == len(line) {continue}
		if i == len(lines) - 1 {lineEnd = ""}
		self.parsingState = self.startState
		self.tag = ""
		self.attrs = make(map[string]string)
		self.remainder = ""
		self.indentCount = 0
		self.closeTag = false
		self.parseLine(line, scope)
		autoclose := self.tagClose()
		switch len(self.tag) {
		case 0:
			output += self.remainder
		default:
			switch len(self.remainder) {
			case 0:
				output += fmt.Sprintf("<%s%s%s>%s", self.tag, self.attrs, autoclose, lineEnd)
			default:
				output += fmt.Sprintf("<%s%s>%s</%s>%s", self.tag, self.attrs, self.remainder, self.tag, lineEnd)
			}
		}
	}
	return
}

func (self *Engine) tagClose() (close string) {
	close = ""
	if self.Options["autoclose"].(bool) || self.closeTag {
		close = " /"
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
