package inputs

import (
	"fmt"
	"strings"

	"github.com/alkemics/goflow"
)

const inputNodeID = "inputs"

type graphRenderer struct {
	goflow.GraphRenderer

	inputs             []goflow.Field
	inputNode          goflow.NodeRenderer
	optionalInputNames []string
}

func (g graphRenderer) Doc() string {
	lines := make([]string, 0, 2)
	doc := g.GraphRenderer.Doc()
	if doc != "" {
		lines = append(lines, doc)
	}
	if len(g.optionalInputNames) > 0 {
		lines = append(lines, fmt.Sprintf(
			"optional: %s",
			strings.Join(g.optionalInputNames, ", "),
		))
	}
	return strings.Join(lines, "\n")
}

func (g graphRenderer) Inputs() []goflow.Field {
	return g.inputs
}

func (g graphRenderer) Nodes() []goflow.NodeRenderer {
	return append(g.GraphRenderer.Nodes(), g.inputNode)
}

type inputNode struct {
	outputs []goflow.Field
}

func (n inputNode) ID() string { return inputNodeID }

func (n inputNode) Previous() []string { return nil }

func (n inputNode) Imports() []goflow.Import { return nil }

func (n inputNode) Doc() string { return "" }

func (n inputNode) Dependencies() []goflow.Field { return nil }

func (n inputNode) Inputs() []goflow.Field {
	return nil
}

func (n inputNode) Outputs() []goflow.Field {
	return n.outputs
}

func (n inputNode) Run(_, outputs []goflow.Field) (string, error) {
	runs := make([]string, len(outputs))
	for i, output := range outputs {
		realOutput := n.outputs[i]
		runs[i] = fmt.Sprintf("%s = %s", output.Name, realOutput.Name)
	}
	return strings.Join(runs, "\n"), nil
}

type graphTypeRenderer struct {
	goflow.GraphRenderer
}

func (g graphTypeRenderer) Inputs() []goflow.Field {
	inputTypes := make(map[string]string)
	for _, n := range g.Nodes() {
		if n.ID() == inputNodeID {
			nodeOutputs := n.Outputs()
			for i, output := range nodeOutputs {
				inputTypes[strings.TrimPrefix(output.Name, "inputs.")] = nodeOutputs[i].Type
			}
		}
	}

	inputs := g.GraphRenderer.Inputs()
	typedInputs := make([]goflow.Field, len(inputs))
	for i, input := range inputs {
		typ := inputTypes[input.Name]
		if typ != "" {
			input.Type = typ
		}
		typedInputs[i] = input
	}

	return typedInputs
}
