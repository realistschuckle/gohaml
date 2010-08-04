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
