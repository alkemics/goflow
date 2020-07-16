package ifs

import (
	"fmt"

	"github.com/alkemics/goflow"
)

type nodeRenderer struct {
	goflow.NodeRenderer

	conditions []condition
}

func (n nodeRenderer) Inputs() []goflow.Field {
	inputs := n.NodeRenderer.Inputs()
	newInputs := make([]goflow.Field, len(inputs)+len(n.conditions))
	copy(newInputs, inputs)

	for i, c := range n.conditions {
		newInputs[len(inputs)+i] = goflow.Field{
			Name: c.name,
			// If we cannot find the type, bool will be used.
			// This can happen if the condition comes from the graph inputs.
			Type: "?bool",
		}
	}
	return newInputs
}

func (n nodeRenderer) Run(inputs, outputs []goflow.Field) (string, error) {
	realInputs := inputs[:len(inputs)-len(n.conditions)]
	conditionInputs := inputs[len(inputs)-len(n.conditions):]

	run, err := n.NodeRenderer.Run(realInputs, outputs)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(` if (%s) {
	%s
}
`, generateConditionals(conditionInputs, n.conditions), run), nil
}
