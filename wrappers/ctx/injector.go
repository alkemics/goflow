package ctx

import "github.com/alkemics/goflow"

type Injector struct {
	goflow.NodeRenderer
}

func (n Injector) Inputs() []goflow.Field {
	inputs := n.NodeRenderer.Inputs()

	// Add context as first input if not present yet
	if len(inputs) > 0 && inputs[0].Type != "context.Context" {
		inputs = append(
			[]goflow.Field{{
				Name: "ctx",
				Type: "context.Context",
			}},
			inputs...,
		)
	}

	return inputs
}

func (n Injector) Run(inputs, outputs []goflow.Field) (string, error) {
	originalInputs := n.NodeRenderer.Inputs()

	if len(originalInputs) > 0 && originalInputs[0].Type != "context.Context" {
		inputs = inputs[1:]
	}

	return n.NodeRenderer.Run(inputs, outputs)
}
