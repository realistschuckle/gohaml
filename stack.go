package gohaml

import (
	"container/list"
	"errors"
)

type stack struct {
	l *list.List
}

func (self *stack) len() (len int) {
	self.lazyInit()
	len = self.l.Len()
	return
}

func (self *stack) push(n Node) {
	self.lazyInit()
	self.l.PushFront(n)
}

func (self *stack) pop() (n Node, err error) {
	self.lazyInit()
	if self.l.Len() == 0 {
		err = errors.New("Can't pop empty list")
		return
	}
	e := self.l.Front()
	self.l.Remove(e)
	n = e.Value.(Node)
	return
}

func (self *stack) lazyInit() {
	if self.l == nil {
		self.l = list.New()
	}
}