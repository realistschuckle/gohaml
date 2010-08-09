package gohaml

import "testing"

type IO struct {
	input string
	expected string
}

var autoCloseTests = []IO{
	IO{"plain ∏ text", "plain ∏ text"},
	IO{"%tag", "<tag />"},
	IO{"%tag tag content", "<tag>tag content</tag>"},
	IO{"%tag.tagClass", "<tag class=\"tagClass\" />"},
	IO{"%tag.tagClass1.tagClass2", "<tag class=\"tagClass1 tagClass2\" />"},
	IO{".tagClass", "<div class=\"tagClass\" />"},
	IO{".tagClass tag content", "<div class=\"tagClass\">tag content</div>"},
	IO{".tagClass1.tagClass2 tag content", "<div class=\"tagClass1 tagClass2\">tag content</div>"},
	IO{"=key1", "value1"},
	IO{"%tag.tagClass= key1", "<tag class=\"tagClass\">value1</tag>"},
	IO{"\\%tag.tagClass= key1", "%tag.tagClass= key1"},
	IO{"\\=key1", "=key1"},
	IO{"%tag#tagId", "<tag id=\"tagId\" />"},
	IO{"#tagId", "<div id=\"tagId\" />"},
	IO{"%tag#tagId.tagClass= key1", "<tag id=\"tagId\" class=\"tagClass\">value1</tag>"},
	IO{"#tagId tag content", "<div id=\"tagId\">tag content</div>"},
	IO{"%tag#tagId= key1", "<tag id=\"tagId\">value1</tag>"},
	IO{"%tag1#tagId1= key1\n%tag2#tagId2= key2", "<tag1 id=\"tagId1\">value1</tag1>\n<tag2 id=\"tagId2\">value2</tag2>"},
	IO{"%tag{:attribute1 => \"value1\", :attribute2 => \"value2\"}", "<tag attribute2=\"value2\" attribute1=\"value1\" />"},
	IO{"%tag{:attribute1 => \"value1\", :attribute2 => \"value2\"} tag content", "<tag attribute2=\"value2\" attribute1=\"value1\">tag content</tag>"},
	IO{"%tag#tagId.tagClass{:id => \"tagId\", :class => \"tagClass\"} tag content", "<tag id=\"tagId tagId\" class=\"tagClass tagClass\">tag content</tag>"},
	IO{"%tag#tagId{:attribute => \"value\"} tag content", "<tag id=\"tagId\" attribute=\"value\">tag content</tag>"},
	IO{"I love <\nHAML!", "I love HAML!"},
	IO{"%a{:href => \"/another/page\"}<\n  %span.button Press me!", "<a href=\"/another/page\"><span class=\"button\">Press me!</span></a>"},
	IO{"%a{:href => \"/another/page\"}<\n  %span.button Press me!\n  %span Me, too!", "<a href=\"/another/page\"><span class=\"button\">Press me!</span>\n<span>Me, too!</span></a>"},
	IO{"%p\n  %a<\n    %span Press me!\n    %span\n      %span Me, too\n    %span And, me!", "<p>\n\t<a><span>Press me!</span>\n\t<span>\n\t\t<span>Me, too</span>\n\t</span>\n\t<span>And, me!</span></a>\n</p>"},
	IO{"I love <\n=lang<\n!", "I love HAML!"},
} 

func TestAutoCloseIO(t *testing.T) {
	for _, io := range autoCloseTests {
		scope := make(map[string]interface{})
		scope["key1"] = "value1"
		scope["key2"] = "value2"
		scope["lang"] = "HAML"
		
		engine := NewEngine(io.input)
		output := engine.Render(scope)
		if output != io.expected {
			t.Errorf("Input    %q\nexpected %q\ngot      %q", io.input, io.expected, output)
		}
	}
}

var noAutoCloseTests = []IO {
	IO{"%tag/", "<tag />"},
	IO{"%tag", "<tag>"},
	IO{"%tag.tagClass", "<tag class=\"tagClass\">"},
	IO{"%tag.tagClass1.tagClass2", "<tag class=\"tagClass1 tagClass2\">"},
	IO{".tagClass", "<div class=\"tagClass\">"},
	IO{"%tag#tagId", "<tag id=\"tagId\">"},
	IO{"#tagId", "<div id=\"tagId\">"},
	IO{"%tag{:attribute1 => \"value1\", :attribute2 => \"value2\"}", "<tag attribute2=\"value2\" attribute1=\"value1\">"},
}

func TestNoAutoCloseIO(t *testing.T) {
	for _, io := range noAutoCloseTests {
		scope := make(map[string]interface{})
		scope["key1"] = "value1"
		scope["key2"] = "value2"
	
		engine := NewEngine(io.input)
		engine.Options["autoclose"] = false
		output := engine.Render(scope)
		if output != io.expected {
			t.Errorf("Input    %q\nexpected %q\ngot      %q", io.input, io.expected, output)
		}
	}
}
