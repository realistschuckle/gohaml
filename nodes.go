package gohaml

type DocTypeNode struct {
	Specification string
}

func (self *DocTypeNode) Accept(c HamlCompiler) {
	c.VisitDocType(self)
}

type TagNode struct {
	Name       string
	ForceClose bool
}

func (self *TagNode) Accept(c HamlCompiler) {
	c.VisitTag(self)
}

type ClassNameNode struct {
	Name string
}

func (self *ClassNameNode) Accept(c HamlCompiler) {
	c.VisitClassName(self)
}
