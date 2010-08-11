package gohaml

import (
	"testing"
)

func TestNewStateDoesNotReturnNil(t* testing.T) {
	state := newState(nil)
	if nil == state {t.Error("new state is nil")}
}

func TestNewStateTransitionsHasValueAndLeftRightIndicesAreZero(t *testing.T) {
	state := newState(nil)
	if nil == state.transitions {t.Error("new state has nil transitions map")}
	if 0 != state.leftIndex {t.Errorf("new state has invalid left index of %d", state.leftIndex)}
	if 0 != state.rightIndex {t.Errorf("new state has invalid left index of %d", state.rightIndex)}
}

func TestInputForUnknownTransitionIncrementsRightIndexAndReturnsSelf(t *testing.T) {
	line := []int("education")
	var scope = make(map[string]interface{})

	state := newState(nil)
	nextState := state.input(0	, line, scope)
	if 1 != state.rightIndex {t.Errorf("expected right index of 1, got %d", state.rightIndex)}
	if state != nextState {t.Error("expected state and next state to be the same")}
}

func TestInputForKnownTransitionDoesNotIncrementRightIndexSetsLeftAndRightIndicesOfDifferentReturnStateToOldRightIndex(t *testing.T) {
	line := []int("education")
	var scope = make(map[string]interface{})
	
	input := 'd'
	from := newState(nil)
	to := newState(nil)
	from.addTransition(func(rune int) bool {return rune == input}, to)
	from.input(0, line, scope)
	newState := from.input(1, line, scope)
	
	if 1 != from.rightIndex {t.Errorf("expected from right index of 1, got %d", from.rightIndex)}
	if 1 != newState.leftIndex {t.Errorf("expected new state left index of 1, got %d", newState.leftIndex)}
	if 2 != newState.rightIndex {t.Errorf("expected new state right index of 2, got %d", newState.rightIndex)}
	if to != newState {t.Error("expected to state and new state to be the same but they weren't")}
}
