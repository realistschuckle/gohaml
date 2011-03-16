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

type io struct {
	input string
	expected string
}

var autoCloseTests = []io{
	io{"plain ∏ text", "plain ∏ text"},
	io{"%tag", "<tag />"},
	io{"%tag tag content", "<tag>tag content</tag>"},
	io{"%tag.tagClass", "<tag class=\"tagClass\" />"},
	io{"%tag.tagClass1.tagClass2", "<tag class=\"tagClass1 tagClass2\" />"},
	io{".tagClass", "<div class=\"tagClass\" />"},
	io{".tagClass tag content", "<div class=\"tagClass\">tag content</div>"},
	io{".tagClass1.tagClass2 tag content", "<div class=\"tagClass1 tagClass2\">tag content</div>"},
	io{"=key1", "value1"},
	io{"%tag.tagClass= key1", "<tag class=\"tagClass\">value1</tag>"},
	io{"\\%tag.tagClass= key1", "%tag.tagClass= key1"},
	io{"\\=key1", "=key1"},
	io{"%tag#tagId", "<tag id=\"tagId\" />"},
	io{"#tagId", "<div id=\"tagId\" />"},
	io{"%tag#tagId.tagClass= key1", "<tag id=\"tagId\" class=\"tagClass\">value1</tag>"},
	io{"#tagId tag content", "<div id=\"tagId\">tag content</div>"},
	io{"%tag#tagId= key1", "<tag id=\"tagId\">value1</tag>"},
	io{"%tag1#tagId1= key1\n%tag2#tagId2= key2", "<tag1 id=\"tagId1\">value1</tag1>\n<tag2 id=\"tagId2\">value2</tag2>"},
	io{"I love <\nHAML!", "I love HAML!"},
	io{"I love <\n=lang<\n!", "I love HAML!"},
	io{"%a{:href => \"/another/page\"}<\n  %span.button Press me!", "<a href=\"/another/page\"><span class=\"button\">Press me!</span></a>"},
	io{"%a{:href => \"/another/page\"}<\n  %span.button Press me!\n  %span Me, too!", "<a href=\"/another/page\"><span class=\"button\">Press me!</span>\n<span>Me, too!</span></a>"},
	io{"%p\n  %a<\n    %span Press me!\n    %span\n      %span Me, too\n    %span And, me!", "<p>\n\t<a><span>Press me!</span>\n\t<span>\n\t\t<span>Me, too</span>\n\t</span>\n\t<span>And, me!</span></a>\n</p>"},
	io{".tagClass{:attribute => key2}", "<div attribute=\"value2\" class=\"tagClass\" />"},
	io{".tagClass{key1 => key2}", "<div value1=\"value2\" class=\"tagClass\" />"},
	io{"#tagId= complexKey.SubKey1", "<div id=\"tagId\">Fortune presents gifts not according to the book.</div>"},
	io{"#tagId= complexKey.SubKey2.SubKey1", "<div id=\"tagId\">That's what I said.</div>"},
	io{"#tagId= complexKey.SubKey2.SubKey2", "<div id=\"tagId\">5</div>"},
	io{"#tagId= complexKey.SubKey2.SubKey3", "<div id=\"tagId\">0.1</div>"},
	io{"#tagId= complexKey.SubKey2.SubKey4.SubKey1", "<div id=\"tagId\">Down deep.</div>"},
	io{"#tagId= complexKey.SubKey2.SubKey4.SubKey2", "<div id=\"tagId\">3</div>"},
	io{"#tagId= complexKey.SubKey2.SubKey4.SubKey3", "<div id=\"tagId\">0.2</div>"},
	io{"#tagId= complexKey.SubKey2.SubKey4.SubKey4", "<div id=\"tagId\" />"},
	io{"=complexKey.SubKey2.SubKey3", "0.1"},
	io{"=complexKey.SubKey3.key", "I got map!"},
	io{"%p= key1", "<p>value1</p>"},
	io{"%tag{:attribute1 => \"value1\", :attribute2 => \"value2\"}", "<tag attribute2=\"value2\" attribute1=\"value1\" />"},
	io{"%tag{:attribute1 => \"value1\", :attribute2 => \"value2\"} tag content", "<tag attribute2=\"value2\" attribute1=\"value1\">tag content</tag>"},
	io{"%tag#tagId.tagClass{:id => \"tagId\", :class => \"tagClass\"} tag content", "<tag id=\"tagId tagId\" class=\"tagClass tagClass\">tag content</tag>"},
	io{"%tag#tagId{:attribute => \"value\"} tag content", "<tag id=\"tagId\" attribute=\"value\">tag content</tag>"},
	io{"%input{:type => \"checkbox\", :checked => true}", "<input type=\"checkbox\" checked=\"checked\" />"},
	io{"%input{:type => \"checkbox\", :checked => false}", "<input type=\"checkbox\" />"},
	io{"%input{:type => \"checkbox\", :checked => outputTrue}", "<input type=\"checkbox\" checked=\"checked\" />"},
	io{"%input{:type => \"checkbox\", cd => outputTrue}", "<input type=\"checkbox\" checked=\"checked\" />"},
	io{"%one\n  %two\n   %three\n", "<one>\n\t<two>\n\t\t<three />\n\t</two>\n</one>"},
	io{"%one\n  %two\n   %three\n      ", "<one>\n\t<two>\n\t\t<three />\n\t</two>\n</one>"},
} 

func TestAutoCloseIO(t *testing.T) {
	for _, io := range autoCloseTests {
		scope := make(map[string]interface{})
		subMap := map[string]interface{} {"key": "I got map!"}
		complexLookup :=  complexLookup{"Fortune presents gifts not according to the book.",
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
			t.Errorf("Input    %q\nexpected %q\ngot      %q", io.input, io.expected, output)
			return
		}
	}
}

var noAutoCloseTests = []io {
	io{"%tag", "<tag>"},
	io{"%tag/", "<tag />"},
	io{"%tag.tagClass", "<tag class=\"tagClass\">"},
	io{"%tag.tagClass1.tagClass2", "<tag class=\"tagClass1 tagClass2\">"},
	io{".tagClass", "<div class=\"tagClass\">"},
	io{"%tag#tagId", "<tag id=\"tagId\">"},
	io{"#tagId", "<div id=\"tagId\">"},
	io{"%tag{:attribute1 => \"value1\", :attribute2 => \"value2\"}", "<tag attribute2=\"value2\" attribute1=\"value1\">"},
}

func TestNoAutoCloseIO(t *testing.T) {
	for _, io := range noAutoCloseTests {
		scope := make(map[string]interface{})
		scope["key1"] = "value1"
		scope["key2"] = "value2"
	
		engine, _ := NewEngine(io.input)
		engine.Autoclose = false
		output := engine.Render(scope)
		if output != io.expected {
			t.Errorf("Input    %q\nexpected %q\ngot      %q", io.input, io.expected, output)
			return
		}
	}
}
