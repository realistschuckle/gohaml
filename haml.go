package gohaml

import (
	"unicode"
	"fmt"
	"strings"
)

type Engine struct {
	Options map[string]interface{}
	Indentation string
	stk *stack
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
	engine = &Engine{make(map[string]interface{}), "\t", newStack(), input, nil, nil, "", make(map[string]string), "", 0, false}
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
		if _, _, hadChildren := self.stk.peek(); i == len(lines) - 1 && !hadChildren {
			lineEnd = ""
		}
		self.parsingState = self.startState
		self.tag = ""
		self.attrs = make(map[string]string)
		self.remainder = ""
		self.indentCount = 0
		self.closeTag = false
		self.parseLine(line, scope)
		
		indent := ""
		for i := 0; i < self.stk.len(); i++ {
			indent += self.Indentation
		}

		if 0 == len(self.tag) {
			output += indent + self.remainder
		} else if len(self.remainder) > 0 {
			if self.stk.len() > 0 {
				output += ">\n"
				self.stk.topHasChildren()
				lineEnd = "\n"
			}
			output += indent + fmt.Sprintf("<%s%s>%s</%s>%s", self.tag, self.attrs, self.remainder, self.tag, lineEnd)
		} else {
			_, topIndent, _ := self.stk.peek()
			if topIndent < self.indentCount {
				if self.stk.len() > 0 {
					output += ">\n"
				}
				self.stk.push(self.tag, self.indentCount)
				output += indent + fmt.Sprintf("<%s%s", self.tag, self.attrs)
			} else if topIndent == self.indentCount {
				autoclose := self.tagClose()
				name, _, hadChildren := self.stk.pop()
				if hadChildren {
					output += indent + fmt.Sprintf("</%s>\n", name)
				} else {
					output += fmt.Sprintf("%s>\n", autoclose)
				}
				indent = ""
				for i := 0; i < self.stk.len(); i++ {
					indent += self.Indentation
				}
				self.stk.push(self.tag, self.indentCount)
				output += indent + fmt.Sprintf("<%s%s", self.tag, self.attrs)
			}
		}
	}
	lineEnd = "\n"
	for depth := self.stk.len(); depth > 0; depth = self.stk.len() {
		autoclose := self.tagClose()
		if 1 == depth {lineEnd = ""}
		indent := ""
		for i := 0; i < depth - 1; i++ {
			indent += self.Indentation
		}
		name, _, hadChildren := self.stk.pop()
		if hadChildren {
			output += indent + fmt.Sprintf("</%s>%s", name, lineEnd)
		} else {
			output += fmt.Sprintf("%s>%s", autoclose, lineEnd)
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
