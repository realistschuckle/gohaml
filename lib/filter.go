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
	Filter(content, indent, oneIndent string) string
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

type FilterFunc func(content, indent, oneIndent string) string

func (fn FilterFunc) Filter(content, indent, oneIndent string) string { return fn(content, indent, oneIndent) }

func cdataHelper(pre, post, content, indent, oneIndent string) string {
	var tail string
	if len(content) > 0 {
		tail = "\n"
	}
	return fmt.Sprintf("%s<![CDATA[%s\n%s%s%s%s]]>%s", pre, post, content, tail, indent, pre, post)
}

func deepenIndent(str, indentation string) string {
	if len(str) == 0 {
		return ""
	}
	return indentation + strings.Replace(str, fmt.Sprintf("\n%s", indentation), fmt.Sprintf("\n%s",strings.Repeat(indentation,2)), -1)
}

func cdata(content, indent, oneIndent string) string { return cdataHelper("", "", content, indent, oneIndent) }
func css(content, indent, oneIndent string) string {
	nextindent := indent + "\t"
	return fmt.Sprintf("<style type=\"text/css\">\n%s%s\n%s</style>",
		nextindent, cdataHelper("/*", "*/", deepenIndent(content, oneIndent), nextindent, oneIndent), indent)
}
func javascript(content, indent, oneIndent string) string {
	nextindent := indent + "\t"
	return fmt.Sprintf("<script type=\"text/javascript\">\n%s%s\n%s</script>",
		nextindent, cdataHelper("//", "", deepenIndent(content, oneIndent), nextindent, oneIndent), indent)
}
