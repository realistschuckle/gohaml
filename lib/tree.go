package gohaml

import (
	"container/vector"
	"fmt"
)

type res struct {
	value string
	needsResolution bool
}

type resPair struct {
	key res
	value res
}

type node struct {
	remainder res
	name string
	attrs vector.Vector
}

type tree struct {
	nodes vector.Vector
}

func newTree() (output *tree) {
	output = &tree{}
	return
}

func (self res) resolve(scope map[string]interface{}) (output string) {
	output = self.value
	return
}

func (self tree) resolve(scope map[string]interface{}) (output string) {
	for _, n := range self.nodes {
		node := n.(*node)
		output += node.resolve(scope)
	}
	return
}

func (self node) resolve(scope map[string]interface{}) (output string) {
	remainder := self.remainder.resolve(scope)
	if self.attrs.Len() > 0 && len(remainder) > 0 {
		if len(self.name) == 0 {self.name = "div"}
		output = fmt.Sprintf("<%s%s>%s</%s>", self.name, self.resolveAttrs(scope), remainder, self.name)
	} else if self.attrs.Len() > 0 {
		if len(self.name) == 0 {self.name = "div"}
		output = fmt.Sprintf("<%s%s />", self.name, self.resolveAttrs(scope))
	} else if len(self.name) > 0 && len(remainder) > 0 {
		output = fmt.Sprintf("<%s>%s</%s>", self.name, remainder, self.name)
	} else if len(self.name) > 0 {
		output = fmt.Sprintf("<%s />", self.name)
	} else {
		output = remainder
	}
	return
}

func (self node) resolveAttrs(scope map[string]interface{}) (output string) {
	attrMap := make(map[string]string)
	for i := 0; i < self.attrs.Len(); i++ {
		resPair := self.attrs.At(i).(*resPair)
		key, value := resPair.key.resolve(scope), resPair.value.resolve(scope)
		if _, ok := attrMap[key]; ok {
			attrMap[key] += " " + value
		} else {
			attrMap[key] = value
		}
	}
	for key, value := range attrMap {
		output += " " + key + "=\"" + value + "\""
	}
	return
}
