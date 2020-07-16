package constants

import "github.com/alkemics/goflow"

type nodeRenderer struct {
	goflow.NodeRenderer

	imports []goflow.Import
	inputs  []goflow.Field
}

func (n nodeRenderer) Imports() []goflow.Import {
	return append(n.NodeRenderer.Imports(), n.imports...)
}

func (n nodeRenderer) Inputs() []goflow.Field {
	return n.inputs
}
