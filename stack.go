package gohaml

type tag struct {
	tagName string
	indent int
	hadChildren bool
}

type entry struct {
	prev *entry
	value tag
}

type stack struct {
	prev *entry
}

func newStack() (s *stack) {
	s = &stack{}
	return
}

func (self *stack) push(tagName string, indent int) {
	entry := &entry{self.prev, tag{tagName, indent, false}}
	if nil != self.prev {self.prev.value.hadChildren = true}
	self.prev = entry
	return
}

func (self *stack) len() (count int) {
	count = 0
	for e := self.prev; nil != e; e = e.prev {
		count++
	}
	return
}

func (self *stack) peek() (tagName string, indent int, hadChildren bool) {
	if nil == self.prev {
		indent = -1
	} else {
		tagName = self.prev.value.tagName
		indent = self.prev.value.indent
		hadChildren = self.prev.value.hadChildren
	}
	return
}

func (self *stack) pop() (tagName string, indent int, hadChildren bool) {
	tagName, indent, hadChildren = self.peek()
	
	if nil != self.prev {self.prev = self.prev.prev}
	return
}

func (self *stack) topHasChildren() {
	if nil != self.prev {
		self.prev.value.hadChildren = true
	}
}
