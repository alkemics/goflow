package types

import "github.com/alkemics/goflow"

type nodeRenderer struct {
	goflow.NodeRenderer

	typedInputs  []goflow.Field
	typedOutputs []goflow.Field
}

func (n nodeRenderer) Inputs() []goflow.Field {
	return n.typedInputs
}

func (n nodeRenderer) Outputs() []goflow.Field {
	return n.typedOutputs
}

type graphRenderer struct {
	goflow.GraphRenderer
	nodes     []goflow.NodeRenderer
	outputMap map[string]string
}

func (g graphRenderer) Nodes() []goflow.NodeRenderer { return g.nodes }

func (g graphRenderer) Outputs() []goflow.Field {
	outputs := g.GraphRenderer.Outputs()
	typedOutputs := make([]goflow.Field, len(outputs))
	for i, f := range outputs {
		if typ, ok := g.outputMap[f.Name]; ok {
			f.Type = typ
		}
		typedOutputs[i] = f
	}
	return typedOutputs
}
