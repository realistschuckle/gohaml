package gohaml

/*
EngineSettings provides the structure of values that the gohaml engine uses
during its parsing and compilation phases.

Has the same fields as the instance attributes of Haml::Engine. You can read
that documentation at http://haml.info/docs/yardoc/Haml/Engine.html
*/
type EngineSettings struct {
	// The string that the engine uses for indentations.
	Indentation string

	// The instance of the HamlCompiler used to compile the templates.
	Compiler HamlCompiler

	// The instance of the HamlParser used to parse the templates.
	Parser HamlParser

	// The options used to create the gohaml engine.
	Options EngineOptions
}
