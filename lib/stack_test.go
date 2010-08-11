package gohaml

import "testing"

func TestLenEqualsZeroForNewStack(t *testing.T) {
	s := newStack()
	if 0 != s.len() {t.Errorf("Expected Len = 0, got %d", s.len())}
}

func TestLenEqualsOneAfterAddingASingleThingToTheStack(t *testing.T) {
	s := newStack()
	s.push("div", 0)
	if 1 != s.len() {t.Errorf("Execpted Len = 1, got %d", s.len())}
}

func TestEmptyIndentCountIsNegativeOne(t *testing.T) {
	s := newStack()
	_, peek, _ := s.peek()
	if -1 != peek {t.Error("Expected empty peek indent to equal -1 but got %d", peek)}
}

func TestPeek(t *testing.T) {
	expectedName, expectedIndent := "tag", 0
	s := newStack()
	s.push(expectedName, expectedIndent)
	
	name, indent, hadChildren := s.peek()
	if expectedName != name {t.Errorf("Expected peek to return %q but got %q", expectedName, name)}
	if expectedIndent != indent {t.Errorf("Expected peek to return %d but got %d", expectedIndent, indent)}
	if hadChildren {t.Errorf("Expected node to not have children but it says it did")}
	if 1 != s.len() {t.Errorf("Expected Len = 1, got %d", s.len())}
	
	expectedName += "1"
	expectedIndent++
	s.push(expectedName, expectedIndent)
	name, indent, hadChildren = s.peek()
	if expectedName != name {t.Errorf("Expected peek to return %q but got %q", expectedName, name)}
	if expectedIndent != indent {t.Errorf("Expected peek to return %d but got %d", expectedIndent, indent)}
	if hadChildren {t.Errorf("Expected node to not have children but it says it did")}
	if 2 != s.len() {t.Errorf("Expected Len = 1, got %d", s.len())}
}

func TestPop(t *testing.T) {
	expectedName, expectedIndent := "div", 0
	s := newStack()
	s.push(expectedName, expectedIndent)
	s.push(expectedName, expectedIndent + 1)

	if 2 != s.len() {t.Errorf("Expected Len = 2, got %d", s.len())}
	
	name, indent, hadChildren := s.pop()
	if expectedName != name {t.Errorf("Expected peek to return %q but got %q", expectedName, name)}
	if expectedIndent + 1 != indent {t.Errorf("Expected peek to return %d but got %d", expectedIndent, indent)}
	if hadChildren {t.Errorf("Expected node to not have children but it says it did")}
	if 1 != s.len() {t.Errorf("Expected Len = 1, got %d", s.len())}

	name, indent, hadChildren = s.pop()
	if expectedName != name {t.Errorf("Expected peek to return %q but got %q", expectedName, name)}
	if expectedIndent != indent {t.Errorf("Expected peek to return %d but got %d", expectedIndent, indent)}
	if !hadChildren {t.Errorf("Expected node to have children but it says it did not")}
	if 0 != s.len() {t.Errorf("Expected Len = 0, got %d", s.len())}
}

func TestSetHasChildren(t *testing.T) {
	s := newStack()
	s.push("tag", 0)
	s.topHasChildren()
	_, _, hadChildren := s.peek()
	if !hadChildren {t.Errorf("Expected node to have children but it says it did not")}
}
