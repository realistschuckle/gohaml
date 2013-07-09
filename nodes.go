package gohaml

type DocTypeNode struct {
	Specification string
}

func (self *DocTypeNode) Accept(c HamlCompiler) {
	c.VisitDocType(self)
}

type TagNode struct {
	Name string
}

func (self *TagNode) Accept(c HamlCompiler) {
	c.VisitTag(self)
}
