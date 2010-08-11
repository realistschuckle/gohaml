package gohaml

type transition struct {
	should func(int) (bool)
	next *state
}

type state struct {
	transitions []*transition
	transitionIndex int
	leftIndex int
	rightIndex int
	exit func(*state, []int, map[string]interface{})
}

func newState(e func(*state, []int, map[string]interface{})) (s *state) {
	return &state{make([]*transition, 10), 0, 0, 0, e}
}

func (self *state) addTransition(f func(int) (bool), s *state) {
	self.transitions[self.transitionIndex] = &transition{f, s}
	self.transitionIndex++
}

func (self *state) input(index int, line []int, scope map[string]interface{}) (s *state) {
	rune := line[index]
	for _, t := range self.transitions {
		if nil == t {break}
		if t.should(rune) {
			if nil != self.exit {self.exit(self, line, scope)}
			s = t.next
			s.leftIndex = index
			s.rightIndex = index + 1
			return
		}
	}
	self.rightIndex = index + 1
	s = self
	return
}
