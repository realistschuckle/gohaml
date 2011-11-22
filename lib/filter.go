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
	if content == "" {
		return fmt.Sprintf("%s<![CDATA[%s\n%s%s]]>%s",pre, post, indent, pre, post)
	}
	return fmt.Sprintf("%s<![CDATA[%s\n%s\n%s%s]]>%s", pre, post, content, indent, pre, post)
}
func cdata(content, indent string) string { return cdataHelper("", "", content, indent) }
func css(content, indent string) string {
	if content == "" {
		return fmt.Sprintf("<style type=\"text/css\">\n%s%s\n%s</style>", indent+"\t", cdataHelper("/*", "*/", content, indent+"\t"),indent)
	}
	return fmt.Sprintf("<style type=\"text/css\">\n%s%s\n%s</style>", indent+"\t", cdataHelper("/*", "*/", "\t" + strings.Replace(content, "\n\t", "\n\t\t", -1), indent+"\t"), indent)
}
func javascript(content, indent string) string {
	if content == "" {
		return fmt.Sprintf("<script type=\"text/javascript\">\n%s%s\n%s</script>", indent+"\t", cdataHelper("//", "", content, indent+"\t"),indent)
	}
	return fmt.Sprintf("<script type=\"text/javascript\">\n%s%s\n%s</script>", indent+"\t", cdataHelper("//", "", "\t" + strings.Replace(content, "\n\t", "\n\t\t", -1), indent+"\t"), indent)
}
