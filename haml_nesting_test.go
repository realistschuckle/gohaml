package gohaml

import "testing"

var nestingTests = []IO{
	IO{"%tag1\n  %tag2", "<tag1>\n	<tag2 />\n<tag1>"},
}

func TestNesting(t *testing.T) {
	for _, io := range nestingTests {
		scope := make(map[string]interface{})
		scope["key1"] = "value1"
		scope["key2"] = "value2"
	
		engine := NewEngine(io.input)
		output := engine.Render(scope)
		if output != io.expected {
			t.Errorf("Input    %q\nexpected %q\ngot      %q", io.input, io.expected, output)
		}
	}
}
