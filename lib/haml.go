// The gohaml package contains a HAML parser similar to the one found at http://www.haml-lang.com.
//
//You can find the specifics about this implementation at http://github.com/realistschuckle/gohaml.
package gohaml

import (
	"unicode"
	"fmt"
	"strings"
	"utf8"
	"reflect"
	"os"
	"io/ioutil"
	"strconv"
)

/*
Engine provides the template interpretation functionality to convert a HAML template into its
corresponding tag-based representation.

Available options are:
  engine.Options["autoclose"] = true|false, default true

The Options field contains the values to modify the way that the engine produces the markup.

The Indentation field contains the string used by the engine to perform indentation.

The IncludeCallback field contains the callback invoked by the gohaml engine to process other files
included through the %include extension.
*/
type Engine struct {
	Options map[string]interface{}
	Indentation string
	IncludeCallback func(string, map[string]interface{}) string
	stk *stack
	input string
	parsingState *state
	startState *state
	tag string
	attrs attrMap
	remainder string
	indentCount int
	closeTag bool
	noNewline bool
	remainderLookup bool
	tree *tree
	currentNode *node
}

// NewEngine returns a new Engine with the given input.
func NewEngine(input string) (engine *Engine) {
	engine = &Engine{make(map[string]interface{}), "\t", loadExternalTemplate, newStack(), input, nil, nil, "", make(map[string]string), "", 0, false, false, false, newTree(), nil}
	engine.makeStates()
	engine.Options["autoclose"] = true
	return
}

// Render uses the provided scope to generate the tag-based representation of the HAML given to the NewFile function.
func (self *Engine) Render(scope map[string]interface{}) (output string) {
	self.tree.indent = self.Indentation
	lines := strings.Split(self.input, "\n", -1)
	if 0 == len(lines) {lines = []string{self.input}}
	for _, line := range lines {
		if 0 == len(line) {continue}
		self.parsingState = self.startState
		self.tag = ""
		self.attrs = make(map[string]string)
		self.remainder = ""
		self.indentCount = 0
		self.closeTag = false
		self.noNewline = false
		self.remainderLookup = false
		self.parseLine(line, scope)
		
		if "include" == self.tag {
			self.remainder = self.IncludeCallback(self.remainder, scope)
			self.tag = ""
		}

		var n *node = nil
		for n = self.currentNode; n != nil && n.indentCount >= self.indentCount; n = self.currentNode.parent {
			self.currentNode = n
		}
		self.currentNode = n
		if nil != self.currentNode {
	 		self.currentNode = self.currentNode.createChild(self.tag, self.remainder, self.indentCount)
		} else {
			self.currentNode = self.tree.createChild(self.tag, self.remainder, self.indentCount)
		}
		if self.remainderLookup {self.currentNode.needsLookup()}
		if !self.tagClose() {self.currentNode.setAutocloseOff()}
		if self.noNewline {self.currentNode.setNoNewline()}
		for key, value := range self.attrs {
			self.currentNode.appendAttr(key, value)
		}
	}
	output = self.tree.String(func(input string) (output string) {return self.lookup(input, scope)})
	return
}

func loadExternalTemplate(path string, scope map[string]interface{}) (output string) {
	in, err := os.Open(path, os.O_RDONLY, 0)
	if nil != err {
		return
	}
	defer func() {
		if nil != in {in.Close()}
	}()
	
	bytes, _ := ioutil.ReadAll(in)
	engine := NewEngine(string(bytes))
	output = engine.Render(scope)
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
		self.remainder = key
		self.remainderLookup = true
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
	
	exitAttributeState := func(s *state, line []int, scope map[string]interface{}) {
		attrPairs := strings.Split(string(line[s.leftIndex + 1:s.rightIndex]), ",", -1)
		for _, attrPair := range attrPairs {
			pair := strings.Split(attrPair, "=>", -1)
			key, value := self.cleanAttr(pair[0], scope), self.cleanAttr(pair[1], scope)
			if value == "true" {value = key}
			if value == "false" {continue}
			self.attrs.appendAttr(key, value)
		}
	}
	
	exitNoNewlineState := func(s *state, line []int, scope map[string]interface{}) {
		self.noNewline = true
	}
	
	exitCodeState := func(s *state, line []int, scope map[string]interface{}) {
		code := string(line[s.leftIndex + 1:s.rightIndex])
		assignment := strings.Split(code, ":=", -1)
		if len(assignment) != 2 {
			// TODO: Report error (Invalid assignment)
		}
		lhs, rhs := self.cleanAttr(assignment[0], nil), self.cleanAttr(assignment[1], scope)
		scope[lhs] = self.convertToType(rhs)
	}
	
	nilFunc := func(s *state, line []int, scope map[string]interface{}) {}
	
 	leadingSpaceState := newState(exitLeadingSpace)
	tagState := newState(exitTagState)
	idState := newState(exitIdState)
	classState := newState(exitClassState)
	keyState := newState(exitKeyState)
	contentState := newState(exitContentState)
	backslashState := newState(exitBackslashState)
	closeTagState := newState(exitCloseTagState)
	attributeState := newState(exitAttributeState)
	endAttributeState := newState(nilFunc)
	noNewlineState := newState(exitNoNewlineState)
	codeState := newState(exitCodeState)
	
	matchTag := func(rune int) bool {return '%' == rune}
	matchId := func(rune int) bool {return '#' == rune}
	matchClass := func(rune int) bool {return '.' == rune}
	matchKey := func(rune int) bool {return '=' == rune}
	matchLeadingContent := func(rune int) bool {return !unicode.IsSpace(rune)}
	matchContent := func(rune int) bool {return unicode.IsSpace(rune)}
	matchBackslashState := func(rune int) bool {return '\\' == rune && 0 == self.indentCount}
	matchCloseTagState := func(rune int) bool {return '/' == rune}
	matchStartAttributeState := func(rune int) bool {return '{' == rune}
	matchEndAttributeState := func(rune int) bool {return '}' == rune}
	matchNoNewlineState := func(rune int) bool {return '<' == rune}
	matchCodeState := func(rune int) bool { return '-' == rune}
	
	leadingSpaceState.addTransition(matchCodeState, codeState)
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
	tagState.addTransition(matchStartAttributeState, attributeState)
	tagState.addTransition(matchNoNewlineState, noNewlineState)
	tagState.addTransition(matchKey, keyState)
	
	attributeState.addTransition(matchEndAttributeState, endAttributeState)
	
	endAttributeState.addTransition(matchContent, contentState)
	endAttributeState.addTransition(matchNoNewlineState, noNewlineState)
	
	keyState.addTransition(matchNoNewlineState, noNewlineState)
	
	idState.addTransition(matchClass, classState)
	idState.addTransition(matchKey, keyState)
	idState.addTransition(matchContent, contentState)
	idState.addTransition(matchStartAttributeState, attributeState)
	
	classState.addTransition(matchClass, classState)
	classState.addTransition(matchKey, keyState)
	classState.addTransition(matchContent, contentState)
	classState.addTransition(matchStartAttributeState, attributeState)
	
	contentState.addTransition(matchNoNewlineState, noNewlineState)
	
	self.parsingState = leadingSpaceState
	self.startState = leadingSpaceState
	
	return
}

func (self *Engine) cleanAttr(attr string, scope map[string]interface{}) (output string) {
	trimFunc := func(rune int) bool {return unicode.IsSpace(rune) || ':' == rune || '"' == rune}
	output = strings.TrimFunc(attr, unicode.IsSpace)
	firstRune, _ := utf8.DecodeRuneInString(output)
	if output == "true" || output == "false" || unicode.IsDigit(firstRune) {
		// Preserve value
	} else if nil != scope && unicode.IsLetter(firstRune) {
		output = fmt.Sprint(scope[output]) // Translate key
	} else {
		output = strings.TrimFunc(output, trimFunc)
	}
	return
}

func (self *Engine) lookup(complexKey string, scope map[string]interface{}) (output string) {
	keyComponents := strings.Split(complexKey, ".", -1)
	if 1 == len(keyComponents) {
		output = fmt.Sprint(scope[complexKey])
	} else {
		var ongoingScope reflect.Value
		for i, keyPart := range keyComponents {
			if 0 == i {
				ongoingScope = reflect.NewValue(scope[keyPart])
			} else {
				var structValue *reflect.StructValue
				
				switch t := ongoingScope.(type) {
				case *reflect.PtrValue:
					structValue = t.Elem().(*reflect.StructValue)
				case *reflect.StructValue:
					structValue = t
				case *reflect.MapValue:
					ongoingScope = t.Elem(reflect.NewValue(keyPart))
					continue
				}
				ongoingScope = structValue.FieldByName(keyPart)
			}
		}
		
		Print:
		switch t := ongoingScope.(type) {
		case *reflect.IntValue:
			output = fmt.Sprint(t.Get())
		case *reflect.StringValue:
			output = fmt.Sprint(t.Get())
		case *reflect.FloatValue:
			output = fmt.Sprint(t.Get())
		case *reflect.InterfaceValue:
			ongoingScope = t.Elem()
			goto Print
		}
	}
	return
}

func (self *Engine) tagClose() (close bool) {
	close = self.Options["autoclose"].(bool) || self.closeTag
	return
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

func (self *Engine) convertToType(input string) (output interface{}) {
	if "true" == input {return true}
	if "false" == input {return false}
	firstRune, _ := utf8.DecodeRuneInString(input)
	if unicode.IsLetter(firstRune) {return input}
	if unicode.IsDigit(firstRune) {
		i, err := strconv.Atoi(input)
		if nil == err {return i}
		f, err := strconv.Atof(input)
		if nil == err {return f}
		// TODO: Report error (cannot convert number)
	}
	// TODO: Report error (unknown parse)
	return
}
