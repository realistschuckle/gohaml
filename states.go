package haml

type state struct {
	transitions map[int]*state
	leftIndex int
	rightIndex int
}

func newState() (s *state) {
	return &state{make(map[int]*state), 0, 0}
}

func (self *state) input(rune int) (s *state) {
	toState, ok := self.transitions[rune]
	if ok {
		toState.leftIndex = self.rightIndex
		toState.rightIndex = toState.leftIndex
		return toState
	}
	self.rightIndex++
	return self
}
