package gohaml

type DocTypeNode struct {
	Specification string
}

func (self *DocTypeNode) Accept(c HamlCompiler) {
	c.VisitDocType(self)
}
