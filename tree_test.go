package gohaml

import "testing"

func TestTreeCreationReturnsNonNilTree(t *testing.T) {
	tree := newTree()
	if nil == tree {t.Error("Expected a non-nil tree")}
}

func TestCreateNodeCreatesNodeInTopLevelNodesAndReturnsNewNode(t *testing.T) {
	tree := newTree()
	node := tree.createChild("name", "remainder", 0)
	if nil == node {t.Error("Expected a non-nil node")}
}

func TestCreateNodeSetsName(t *testing.T) {
	expectedName := "node name"
	node := newTree().createChild(expectedName, "remainder", 0)
	if node.name != expectedName {t.Errorf("Expected node name %q but got %q", expectedName, node.name)}
}

func TestCreateNodeCreatesNewAttributeMapForNode(t *testing.T) {
	attrName, attrValue := "class", "woot"
	node := newTree().createChild("name", "remainder", 0)
	node.attrs[attrName] = attrValue
	if attrValue != node.attrs[attrName] {t.Error("Setting an attribute failed")}
}

func TestCreateNodeSetsRemainder(t *testing.T) {
	remainder := "some stuff to print"
	node := newTree().createChild("name", remainder, 0)
	if remainder != node.remainder {t.Errorf("Expected remainder %q but got %q", remainder, node.remainder)}
}

func TestChildNodeCreationShouldActJustLikeItDoesOnTree(t *testing.T) {
	childName, childRemainder := "child node", "child remainder"
	childAttrName, childAttrValue := "child attr", "child attr value"
	node := newTree().createChild("name", "remainder", 0)
	node = node.createChild(childName, childRemainder, 0)

	if nil == node {t.Error("Expected a non-nil child node"); return}
	if childName != node.name {t.Errorf("Expected node name %q but got %q", childName, node.name)}
	if childRemainder != node.remainder {t.Errorf("Expected remainder %q but got %q", childRemainder, node.remainder)}

	node.attrs[childAttrName] = childAttrValue
	if childAttrValue != node.attrs[childAttrName] {t.Error("Setting an attribute failed")}
}

func TestParentOfChildNodeIsSet(t *testing.T) {
	node := newTree().createChild("", "", 0)
	childNode := node.createChild("", "", 0)
	if node != childNode.parent {t.Errorf("Parent node of child %s not set to the real parent %s", childNode.parent, node)}
}

func TestChildNodeOfTreeHasNilParent(t *testing.T) {
	node := newTree().createChild("", "", 0)
	if !node.topLevel() {t.Error("Expected top-level node to report itself as such")}
}

func TestChildNodeOfNodeIsNotTopLevel(t *testing.T) {
	node := newTree().createChild("", "", 0).createChild("", "", 0)
	if node.topLevel() {t.Error("Expected child node to not be top level but it says it is")}
}

func TopLevelChildAppearsInTreeNodeChildrenCollection(t *testing.T) {
	tree := newTree()
	node := tree.createChild("", "", 0)
	if node != tree.childAt(0) {t.Error("New child not at beginning of tree children list")}
}
