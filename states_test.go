package gohaml

import (
	"testing"
)

func TestNewStateDoesNotReturnNil(t* testing.T) {
	state := newState()
	if nil == state {t.Error("new state is nil")}
}

func TestNewStateTransitionsHasValueAndLeftRightIndicesAreZero(t *testing.T) {
	state := newState()
	if nil == state.transitions {t.Error("new state has nil transitions map")}
	if 0 != state.leftIndex {t.Errorf("new state has invalid left index of %d", state.leftIndex)}
	if 0 != state.rightIndex {t.Errorf("new state has invalid left index of %d", state.rightIndex)}
}

func TestInputForUnknownTransitionIncrementsRightIndexAndReturnsSelf(t *testing.T) {
	state := newState()
	nextState := state.input('d')
	if 1 != state.rightIndex {t.Errorf("expected right index of 1, got %d", state.rightIndex)}
	if state != nextState {t.Error("expected state and next state to be the same")}
}

func TestInputForKnownTransitionDoesNotIncrementRightIndexSetsLeftAndRightIndicesOfDifferentReturnStateToOldRightIndex(t *testing.T) {
	input := 'd'
	from := newState()
	to := newState()
	from.transitions[input] = to
	from.input('e')
	newState := from.input(input)
	
	if 1 != from.rightIndex {t.Errorf("expected from right index of 1, got %d", from.rightIndex)}
	if 1 != newState.leftIndex {t.Errorf("expected new state left index of 1, got %d", newState.leftIndex)}
	if 1 != newState.rightIndex {t.Errorf("expected new state right index of 1, got %d", newState.rightIndex)}
	if to != newState {t.Error("expected to state and new state to be the same but they weren't")}
}
