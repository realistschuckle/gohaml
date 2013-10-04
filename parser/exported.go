package parser

type ParsedDocument struct {
	// Indentation
}

type HamlParser interface {
	Parse(string) (ParsedDocument, error)
}
