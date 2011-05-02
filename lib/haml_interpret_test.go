package gohaml

import (
	"testing"
	"strings"
	"fmt"
)

type assignment struct {
	name string
	rhs string
	value interface{}
}

func TestForSliceRangeConstruct(t *testing.T) {
	scope := make(map[string]interface{})
	scope["looper"] = []int{4, -128, 38, 99, 1}
	
	expected := "<p>\n" +
				"	<span>0</span><span>4</span>\n" +
				"	<span>1</span><span>-128</span>\n" +
				"	<span>2</span><span>38</span>\n" +
				"	<span>3</span><span>99</span>\n" +
				"	<span>4</span><span>1</span>\n" +
				"</p>"
	input := "%p\n  - for i, v := range looper\n    %span= i<\n    %span= v"
	engine, _ := NewEngine(input)
	output := engine.Render(scope)
	
	if output != expected {
		t.Errorf("Expected\n%s\nbut got\n%s\n", expected, output)
	}
}

func TestForArrayRangeConstruct(t *testing.T) {
	scope := make(map[string]interface{})
	scope["looper"] = [5]int{4, -128, 38, 99, 1}
	
	expected := "<p>\n" +
				"	<span>0</span><span>4</span>\n" +
				"	<span>1</span><span>-128</span>\n" +
				"	<span>2</span><span>38</span>\n" +
				"	<span>3</span><span>99</span>\n" +
				"	<span>4</span><span>1</span>\n" +
				"</p>"
	input := "%p\n  - for i, v := range looper\n    %span= i<\n    %span= v"
	engine, _ := NewEngine(input)
	output := engine.Render(scope)
	
	if output != expected {
		t.Errorf("Expected\n%s\nbut got\n%s\n", expected, output)
	}
}

func TestForMapRangeConstruct(t *testing.T) {
	scope := make(map[string]interface{})
	intmap := map[int]int {
		0:4,
		1:-128,
		2:38,
		3:99,
		4:1,
	}
	scope["looper"] = intmap

	input := "%p\n  - for i, v := range looper\n    %span= i<\n    %span= v"
	engine, _ := NewEngine(input)
	output := engine.Render(scope)

	for k, v := range intmap {
		expect := fmt.Sprintf("<span>%d</span><span>%d</span>", k, v)
		if !strings.Contains(output, expect) {
			t.Errorf("Execpted\n%s\nbut got\n%s\n", expect, output)
		}
	}
}

var assignmentInputs = []assignment {
	assignment{"localString", "\"string\"", "string"},
	assignment{"localInt", "3", 3},
	assignment{"localFloat", "3.14", 3.14},
	assignment{"localLookup", "commonKey", "commonValue"},
	assignment{"localLookup", "akey.subkey", "subkeyvalue"},
}

type subkey struct {
	subkey string
}

func TestAssignments(t *testing.T) {
	for _, assignment := range assignmentInputs {
		scope := make(map[string]interface{})
		scope["commonKey"] = "commonValue"
		scope["akey"] = &subkey{"subkeyvalue"}

		for _, input := range generateAssignments(assignment) {
			engine, _ := NewEngine(input)
			output := engine.Render(scope)

			if _, ok := scope[assignment.name]; !ok {
				s := fmt.Sprint(scope)
				t.Errorf("Input %q\nMap   %s", input, s)
				return
			}

			if scope[assignment.name] != assignment.value {
				s := fmt.Sprint(scope)
				t.Errorf("Input %q\nMap   %s", input, s)
				return
			}
			
			if len(output) > 0 {
				t.Errorf("Expected no output, but got %q", output)
				return
			}
		}
	}
}

func generateAssignments(assignment assignment) (assignments []string) {
	return []string {
		fmt.Sprintf("- %s := %s", assignment.name, assignment.rhs),
		fmt.Sprintf("-%s := %s", assignment.name, assignment.rhs),
		fmt.Sprintf("- %s:= %s", assignment.name, assignment.rhs),
		fmt.Sprintf("- %s :=%s", assignment.name, assignment.rhs),
		fmt.Sprintf("-%s:= %s", assignment.name, assignment.rhs),
		fmt.Sprintf("-%s :=%s", assignment.name, assignment.rhs),
		fmt.Sprintf("-%s:=%s", assignment.name, assignment.rhs),
		fmt.Sprintf("  - %s := %s  ", assignment.name, assignment.rhs),
	}
}
