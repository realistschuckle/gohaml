package gohaml

/*  Filename:    filter.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Mon Nov 21 08:55:45 PST 2011
 *  Description: 
 */

import (
	"strings"
	"fmt"
)

type Filter interface {
	// Content is indented and ends with "\n".
	// Output should be indented and end with "\n".
	Filter(content, indent string) string
}

type FilterMap map[string]Filter

func (filters FilterMap) Copy() FilterMap {
	cp := make(FilterMap, len(filters))
	for k, v := range filters {
		cp[k] = v
	}
	return cp
}

var defaultFilterMap = FilterMap{
	"cdata":      FilterFunc(cdata),
	"css":        FilterFunc(css),
	"javascript": FilterFunc(javascript),
}

type FilterFunc func(content, indent string) string

func (fn FilterFunc) Filter(content, indent string) string { return fn(content, indent) }

func cdataHelper(pre, post, content, indent string) string {
	var tail string
	if len(content) > 0 {
		tail = "\n"
	}
	return fmt.Sprintf("%s<![CDATA[%s\n%s%s%s%s]]>%s", pre, post, content, tail, indent, pre, post)
}

func deepenIndent(str string) string {
	if len(str) == 0 {
		return ""
	}
	return "\t" + strings.Replace(str, "\n\t", "\n\t\t", -1)
}

func cdata(content, indent string) string { return cdataHelper("", "", content, indent) }
func css(content, indent string) string {
	nextindent := indent + "\t"
	return fmt.Sprintf("<style type=\"text/css\">\n%s%s\n%s</style>",
		nextindent, cdataHelper("/*", "*/", deepenIndent(content), nextindent), indent)
}
func javascript(content, indent string) string {
	nextindent := indent + "\t"
	return fmt.Sprintf("<script type=\"text/javascript\">\n%s%s\n%s</script>",
		nextindent, cdataHelper("//", "", deepenIndent(content), nextindent), indent)
}
