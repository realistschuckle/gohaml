package gohaml

import "testing"

type simpleLookup struct {
	SubKey1 string
	SubKey2 int
	SubKey3 float64
	SubKey4 *simpleLookup
}

type complexLookup struct {
	SubKey1 string
	SubKey2 simpleLookup
	SubKey3 map[string]interface{}
}

type testcase struct {
	input    string
	expected string
}

var autoCloseTests = []testcase{
	testcase{"plain ∏ text", "plain ∏ text"},
	testcase{"%tag", "<tag />"},
	testcase{"%tag tag content", "<tag>tag content</tag>"},
	testcase{"%tag.tagClass", "<tag class=\"tagClass\" />"},
	testcase{"%tag.tagClass1.tagClass2", "<tag class=\"tagClass1 tagClass2\" />"},
	testcase{".tagClass", "<div class=\"tagClass\" />"},
	testcase{".tagClass tag content", "<div class=\"tagClass\">tag content</div>"},
	testcase{".tagClass1.tagClass2 tag content", "<div class=\"tagClass1 tagClass2\">tag content</div>"},
	testcase{"=key1", "value1"},
	testcase{"%tag.tagClass= key1", "<tag class=\"tagClass\">value1</tag>"},
	testcase{"\\%tag.tagClass= key1", "%tag.tagClass= key1"},
	testcase{"\\=key1", "=key1"},
	testcase{"%tag#tagId", "<tag id=\"tagId\" />"},
	testcase{"#tagId", "<div id=\"tagId\" />"},
	testcase{"%tag#tagId.tagClass= key1", "<tag id=\"tagId\" class=\"tagClass\">value1</tag>"},
	testcase{"#tagId tag content", "<div id=\"tagId\">tag content</div>"},
	testcase{"%tag#tagId= key1", "<tag id=\"tagId\">value1</tag>"},
	testcase{"%tag1#tagId1= key1\n%tag2#tagId2= key2", "<tag1 id=\"tagId1\">value1</tag1>\n<tag2 id=\"tagId2\">value2</tag2>"},
	testcase{"I love <\nHAML!", "I love HAML!"},
	testcase{"I love <\n=lang<\n!", "I love HAML!"},
	testcase{"%a{:href => \"/another/page\"}<\n  %span.button Press me!", "<a href=\"/another/page\"><span class=\"button\">Press me!</span></a>"},
	testcase{"%a{:href => \"/another/page\"}<\n  %span.button Press me!\n  %span Me, too!", "<a href=\"/another/page\"><span class=\"button\">Press me!</span>\n<span>Me, too!</span></a>"},
	testcase{"%p\n  %a<\n    %span Press me!\n    %span\n      %span Me, too\n    %span And, me!", "<p>\n\t<a><span>Press me!</span>\n\t<span>\n\t\t<span>Me, too</span>\n\t</span>\n\t<span>And, me!</span></a>\n</p>"},
	testcase{".tagClass{:attribute => key2}", "<div class=\"tagClass\" attribute=\"value2\" />"},
	testcase{".tagClass{key1 => key2}", "<div class=\"tagClass\" value1=\"value2\" />"},
	testcase{"#tagId= complexKey.SubKey1", "<div id=\"tagId\">Fortune presents gifts not according to the book.</div>"},
	testcase{"#tagId= complexKey.SubKey2.SubKey1", "<div id=\"tagId\">That's what I said.</div>"},
	testcase{"#tagId= complexKey.SubKey2.SubKey2", "<div id=\"tagId\">5</div>"},
	testcase{"#tagId= complexKey.SubKey2.SubKey3", "<div id=\"tagId\">0.1</div>"},
	testcase{"#tagId= complexKey.SubKey2.SubKey4.SubKey1", "<div id=\"tagId\">Down deep.</div>"},
	testcase{"#tagId= complexKey.SubKey2.SubKey4.SubKey2", "<div id=\"tagId\">3</div>"},
	testcase{"#tagId= complexKey.SubKey2.SubKey4.SubKey3", "<div id=\"tagId\">0.2</div>"},
	testcase{"#tagId= complexKey.SubKey2.SubKey4.SubKey4", "<div id=\"tagId\" />"},
	testcase{"=complexKey.SubKey2.SubKey3", "0.1"},
	testcase{"=complexKey.SubKey3.key", "I got map!"},
	testcase{"%p= key1", "<p>value1</p>"},
	testcase{"%tag{:attribute1 => \"value1\", :attribute2 => \"value2\"}", "<tag attribute1=\"value1\" attribute2=\"value2\" />"},
	testcase{"%tag{:attribute1 => \"value1\", :attribute2 => \"value2\"} tag content", "<tag attribute1=\"value1\" attribute2=\"value2\">tag content</tag>"},
	testcase{"%tag#tagId.tagClass{:id => \"tagId\", :class => \"tagClass\"} tag content", "<tag id=\"tagId tagId\" class=\"tagClass tagClass\">tag content</tag>"},
	testcase{"%tag#tagId{:attribute => \"value\"} tag content", "<tag id=\"tagId\" attribute=\"value\">tag content</tag>"},
	testcase{"%input{:type => \"checkbox\", :checked => true}", "<input type=\"checkbox\" checked=\"checked\" />"},
	testcase{"%input{:type => \"checkbox\", :checked => false}", "<input type=\"checkbox\" />"},
	testcase{"%input{:type => \"checkbox\", :checked => outputTrue}", "<input type=\"checkbox\" checked=\"checked\" />"},
	testcase{"%input{:type => \"checkbox\", cd => outputTrue}", "<input type=\"checkbox\" checked=\"checked\" />"},
	testcase{"%one\n  %two\n   %three\n", "<one>\n\t<two>\n\t\t<three />\n\t</two>\n</one>"},
	testcase{"%one\n  %two\n   %three\n      ", "<one>\n\t<two>\n\t\t<three />\n\t</two>\n</one>"},
	testcase{"!!!", "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">"},
	testcase{"!!! Strict", "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\">"},
	testcase{"!!! Frameset", "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Frameset//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-frameset.dtd\">"},
	testcase{"!!! 5", "<!DOCTYPE html>"},
	testcase{"!!! 1.1", "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.1//EN\" \"http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd\">"},
	testcase{"!!! Basic", "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML Basic 1.1//EN\" \"http://www.w3.org/TR/xhtml-basic/xhtml-basic11.dtd\">"},
	testcase{"!!! Mobile", "<!DOCTYPE html PUBLIC \"-//WAPFORUM//DTD XHTML Mobile 1.2//EN\" \"http://www.openmobilealliance.org/tech/DTD/xhtml-mobile12.dtd\">"},
	testcase{"!!! RDFa", "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML+RDFa 1.0//EN\" \"http://www.w3.org/MarkUp/DTD/xhtml-rdfa-1.dtd\">"},
}

func TestAutoCloseIO(t *testing.T) {
	for i, io := range autoCloseTests {
		scope := make(map[string]interface{})
		subMap := map[string]interface{}{"key": "I got map!"}
		complexLookup := complexLookup{"Fortune presents gifts not according to the book.",
			simpleLookup{"That's what I said.", 5, .1,
				&simpleLookup{"Down deep.", 3, .2, nil}},
			nil}
		complexLookup.SubKey3 = subMap
		scope["complexKey"] = complexLookup
		scope["key1"] = "value1"
		scope["key2"] = "value2"
		scope["lang"] = "HAML"
		scope["outputFalse"] = "false"
		scope["outputTrue"] = "true"
		scope["cd"] = "checked"

		engine, _ := NewEngine(io.input)
		output := engine.Render(scope)
		if output != io.expected {
			t.Errorf("(%d) Input    %q\nexpected %q\ngot      %q", i, io.input, io.expected, output)
			return
		}
	}
}

var noAutoCloseTests = []testcase{
	testcase{"%tag", "<tag>"},
	testcase{"%tag/", "<tag />"},
	testcase{"%tag.tagClass", "<tag class=\"tagClass\">"},
	testcase{"%tag.tagClass1.tagClass2", "<tag class=\"tagClass1 tagClass2\">"},
	testcase{".tagClass", "<div class=\"tagClass\">"},
	testcase{"%tag#tagId", "<tag id=\"tagId\">"},
	testcase{"#tagId", "<div id=\"tagId\">"},
	testcase{"%tag{:attribute1 => \"value1\", :attribute2 => \"value2\"}", "<tag attribute1=\"value1\" attribute2=\"value2\">"},
}

func TestNoAutoCloseIO(t *testing.T) {
	for i, io := range noAutoCloseTests {
		scope := make(map[string]interface{})
		scope["key1"] = "value1"
		scope["key2"] = "value2"

		engine, _ := NewEngine(io.input)
		engine.Autoclose = false
		output := engine.Render(scope)
		if output != io.expected {
			t.Errorf("(%d)Input    %q\nexpected %q\ngot      %q", i, io.input, io.expected, output)
			return
		}
	}
}
