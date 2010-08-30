package gohaml

import (
	"testing"
	"fmt"
)

type assignment struct {
	name string
	rhs string
	value interface{}
}

var assignmentInputs = []assignment {
	assignment{"localString", "\"string\"", "string"},
	assignment{"localInt", "3", 3},
	assignment{"localFloat", "3.14", 3.14},
	assignment{"localLookup", "commonKey", "commonValue"},
}

func TestAssignments(t *testing.T) {
	for _, assignment := range assignmentInputs {
		scope := make(map[string]interface{})
		scope["commonKey"] = "commonValue"

		for _, input := range generateAssignments(assignment) {
			engine := NewEngine(input)
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
