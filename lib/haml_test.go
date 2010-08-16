package gohaml

import (
	"testing"
	"os"
)

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

func TestCustomIndentPropertyAffectOnOutput(t *testing.T) {
	engine := NewEngine("%tag\n  %subTag")
	engine.Indentation = "    "
	output := engine.Render(nil)
	if output != "<tag>\n    <subTag />\n</tag>" {
		t.Errorf("Expected custom-indented tags but got %q", output)
	}
}

func TestIncludeCallbackWorks(t *testing.T) {
	callbackCalledCorrectly := false
	templatePath := "my/include.haml"
	templateScope := make(map[string]interface{})
	engine := NewEngine("%include " + templatePath)
	engine.IncludeCallback = func(include string, scope map[string]interface{}) (output string) {
		callbackCalledCorrectly = include == templatePath && scope == templateScope
		return
	}
	engine.Render(templateScope)
	
	if !callbackCalledCorrectly {t.Error("Expected callback invocation with path and scope")}
}

func TestDefaultIncludeWorks(t *testing.T) {
	defer func() {
		os.RemoveAll("includes")
	}()
	err := writeFile("includes", "test.haml", "%p= key1")
	if nil != err {t.Error(err); return}

	scope := map[string]interface{} {"key1" : "Include extension works!"}
	engine := NewEngine("%include includes/test.haml")
	output := engine.Render(scope)
	expectedOutput := "<p>Include extension works!</p>"
	if output != expectedOutput {
		t.Errorf("Expected %q\ngot      %q", expectedOutput, output)
	}
}

func writeFile(dir string, filename string, content string) (err os.Error) {
	includePath := dir + "/" + filename;
	{
		err = os.Mkdir(dir, 0777)
		if nil != err {return}
	}
	{
		var out *os.File
		out, err = os.Open(includePath, os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0666)
		if nil != err { return }
		out.WriteString(content)
		out.Close()
	}
	return
}
