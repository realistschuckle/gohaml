package gohaml

import "testing"

func TestNewEngine(t *testing.T) {
	engine := NewEngine("world")
	if nil == engine {t.Error("Engine is nil.")}
}

func TestNewEngineHasOptionsMapCreated(t *testing.T) {
	engine := NewEngine("world")
	if nil == engine.Options {t.Error("Options map is nil.")}
}

func TestIndentCount(t *testing.T) {
	checkIndent(1, "%tag\n %tag", t)
	checkIndent(4, "%tag\n  %tag\n    %tag", t)
	checkIndent(2, "%tag\n  %tag", t)
	checkIndent(0, "%tag", t)
}

func checkIndent(expectedCount int, input string, t *testing.T) {
	engine := NewEngine(input)
	engine.Render(nil)
	if expectedCount != engine.indentCount {
		t.Errorf("Expected indent count of %d but got %d.", expectedCount, engine.indentCount)
	}
}
