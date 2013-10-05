package parser

type ParsedDocument struct {
	// Indentation
}

type HamlParser interface {
	Parse(string) (ParsedDocument, error)
}

type DefaultParser struct {

}

func (self *DefaultParser) Parse(string) (d ParsedDocument, e error) {
	return
}
