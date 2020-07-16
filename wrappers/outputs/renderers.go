package outputs

import (
	"fmt"
	"regexp"

	"github.com/alkemics/goflow"
	"github.com/alkemics/goflow/wrappers/bind"
)

type outputNode struct {
	id     string
	inputs []goflow.Field
	output goflow.Field
}

func (n outputNode) ID() string {
	return n.id
}

func (n outputNode) Doc() string { return "" }

func (n outputNode) Previous() []string {
	return nil
}

func (n outputNode) Imports() []goflow.Import {
	return nil
}

func (n outputNode) Dependencies() []goflow.Field { return nil }

func (n outputNode) Inputs() []goflow.Field {
	return n.inputs
}

func (n outputNode) Outputs() []goflow.Field {
	return []goflow.Field{n.output}
}

func (n outputNode) Run(inputs, outputs []goflow.Field) (string, error) {
	run, err := bind.Run(inputs, outputs)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`
%s
%s = %s`,
		run,
		n.output.Name,
		outputs[0].Name,
	), nil
}

type graphRenderer struct {
	goflow.GraphRenderer

	outputs        []goflow.Field
	outputBuilders []goflow.NodeRenderer
}

func (g graphRenderer) Outputs() []goflow.Field {
	return g.outputs
}

func (g graphRenderer) Nodes() []goflow.NodeRenderer {
	return append(g.GraphRenderer.Nodes(), g.outputBuilders...)
}

var outputBuilderRegex = regexp.MustCompile("^__output_(.*)_builder")

type nameGraphWrapper struct {
	goflow.GraphRenderer
	outputs []goflow.Field
}

func (g nameGraphWrapper) Outputs() []goflow.Field { return g.outputs }
