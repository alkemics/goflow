package goflow

type GraphRenderer interface {
	Name() string
	Pkg() string
	Imports() []Import
	Doc() string
	Dependencies() []Field
	Inputs() []Field
	Outputs() []Field
	Nodes() []NodeRenderer
}

type NodeRenderer interface {
	ID() string
	Previous() []string
	Imports() []Import
	Doc() string
	Dependencies() []Field
	Inputs() []Field
	Outputs() []Field
	Run(inputs, outputs []Field) (string, error)
}

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func FieldNames(fields []Field) []string {
	fieldNames := make([]string, len(fields))
	for i, f := range fields {
		fieldNames[i] = f.Name
	}
	return fieldNames
}

type Import struct {
	Pkg string
	Dir string
}
