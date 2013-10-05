package parser

import (
	"errors"
	"io"
	"unicode/utf8"
)

type scanner struct {
	reader io.RuneReader
	memory [8]rune
	memptr int
	unread int
}

func (self *scanner) ReadRune() (r rune, size int, err error) {
	if self.unread > 0 {
		r = self.memory[(self.memptr-self.unread)%8]
		size = utf8.RuneLen(r)
		self.unread -= 1
	} else {
		r, size, err = self.reader.ReadRune()
		self.memory[self.memptr%8] = r
		self.memptr += 1
	}
	return
}

func (self *scanner) UnreadRune() (err error) {
	switch {
	case self.unread > 7:
		err = errors.New("scanner cannot unread more than 8 runes.")
	case self.unread == self.memptr:
		err = errors.New("scanner cannot unread more than read")
	default:
		self.unread += 1
	}
	return
}
