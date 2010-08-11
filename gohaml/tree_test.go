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

func TestStringRepresentationOfAnEmptyTreeIsAnEmptyString(t *testing.T) {
	ts := newTree().String()
	if ts != "" {t.Errorf("Expected empty tree string representation to be empty but got %q", ts)}
}

func TestStringRepOfNodeWithNoNameIsOnlyRemainder(t *testing.T) {
	expectedValue := "remainder"
	tree := newTree()
	tree.createChild("", expectedValue, 0)
	ts := tree.String()
	if expectedValue != ts {t.Errorf("Expected tree to be %q but got %q", expectedValue, ts)}
}

func TestStringRepOfNodeWithTagIsTagNameOnly(t *testing.T) {
	expectedValue := "<tag1 />"
	tree := newTree()
	tree.createChild("tag1", "", 0)
	ts := tree.String()
	if expectedValue != ts {t.Errorf("Expected tree to be %q but got %q", expectedValue, ts)}
}

func TestStringRepOfNodeWithTagAndRemainderIsTagAndRemainder(t *testing.T) {
	expectedValue := "<tag1>tag content</tag1>"
	tree := newTree()
	tree.createChild("tag1", "tag content", 0)
	ts := tree.String()
	if expectedValue != ts {t.Errorf("Expected tree to be %q but got %q", expectedValue, ts)}
}

func TestStringRepOfNodeWithAttributes(t *testing.T) {
	expectedValue := "<tag1 id=\"tagId\" class=\"tagClass1 tagClass2\">tag content</tag1>"
	tree := newTree()
	node := tree.createChild("tag1", "tag content", 0)
	node.appendAttr("id", "tagId")
	node.appendAttr("class", "tagClass1")
	node.appendAttr("class", "tagClass2")
	ts := tree.String()
	if expectedValue != ts {t.Errorf("Expected tree to be %q but got %q", expectedValue, ts)}
}

func TestStringRepOfMutlipleTopLevelNodes(t *testing.T) {
	expectedValue := "<tag1>tag content 1</tag1>\n<tag2>tag content 2</tag2>"
	tree := newTree()
	tree.createChild("tag1", "tag content 1", 0)
	tree.createChild("tag2", "tag content 2", 0)
	ts := tree.String()
	if expectedValue != ts {t.Errorf("Expected tree to be %q but got %q", expectedValue, ts)}
}

func TestNestedTreeStringRep(t *testing.T) {
	expectedValue := "<root>\n\t<child1>child 1 content</child1>\n\t<child2 />\n\tPlain text\n</root>"
	tree := newTree()
	root := tree.createChild("root", "", 0)
	root.createChild("child1", "child 1 content", 0)
	root.createChild("child2", "", 0)
	root.createChild("", "Plain text", 0)
	ts := tree.String()
	if expectedValue != ts {t.Errorf("Expected tree to be %q\nbut got             %q", expectedValue, ts)}
}

func TestDeepNestedTreeStringRep(t *testing.T) {
	expectedValue := "<root>\n\t<child1>\n\t\t<child2>\n\t\t\t<child3 />\n\t\t</child2>\n\t</child1>\n</root>"
	tree := newTree()
	root := tree.createChild("root", "", 0)
	node := root.createChild("child1", "", 0)
	node = node.createChild("child2", "", 0)
	node.createChild("child3", "", 0)
	ts := tree.String()
	if expectedValue != ts {t.Errorf("Expected tree to be %q\nbut got             %q", expectedValue, ts)}
}

func TestTurnOffCloseTag(t *testing.T) {
	expectedValue := "<tag>"
	tree := newTree()
	root := tree.createChild("tag", "", 0)
	root.setAutocloseOff()
	
	ts := tree.String()
	if expectedValue != ts {t.Errorf("Expected tree to be %q\nbut got             %q", expectedValue, ts)}
}
