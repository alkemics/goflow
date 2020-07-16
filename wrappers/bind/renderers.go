package bind

import (
	"fmt"
	"strings"

	"github.com/alkemics/goflow"
)

type nodeRenderer struct {
	goflow.NodeRenderer

	bindMap          map[string]goflow.Field
	absentInputNames map[string]struct{}
}

func (n nodeRenderer) Inputs() []goflow.Field {
	boundInputs := make([]goflow.Field, len(n.NodeRenderer.Inputs()))
	for i, input := range n.NodeRenderer.Inputs() {
		f := input
		if bf, ok := n.bindMap[f.Name]; ok {
			f = bf
		}
		boundInputs[i] = f
	}
	return boundInputs
}

func (n nodeRenderer) Run(inputs, outputs []goflow.Field) (string, error) {
	statements := make([]string, len(n.absentInputNames)+1)
	i := 0
	for _, input := range inputs {
		if _, ok := n.absentInputNames[input.Name]; !ok {
			continue
		}
		statements[i] = fmt.Sprintf("var %s %s", input.Name, input.Type)
		i++
	}

	stmts, err := n.NodeRenderer.Run(inputs, outputs)
	if err != nil {
		return "", err
	}

	statements[len(n.absentInputNames)] = stmts
	return strings.Join(statements, "\n"), nil
}

type bindingNode struct {
	id     string
	inputs []goflow.Field
	output goflow.Field
}

func (n bindingNode) ID() string {
	return n.id
}

func (n bindingNode) Doc() string { return "" }

func (n bindingNode) Previous() []string {
	return nil
}

func (n bindingNode) Imports() []goflow.Import {
	return nil
}

func (n bindingNode) Dependencies() []goflow.Field {
	return nil
}

func (n bindingNode) Inputs() []goflow.Field {
	return n.inputs
}

func (n bindingNode) Outputs() []goflow.Field {
	return []goflow.Field{n.output}
}

func (n bindingNode) Run(inputs, outputs []goflow.Field) (string, error) {
	return Run(inputs, outputs)
}

type graphRenderer struct {
	goflow.GraphRenderer

	wrappedNodes []goflow.NodeRenderer
	bindingNodes []bindingNode
}

func (g graphRenderer) Nodes() []goflow.NodeRenderer {
	nodes := make([]goflow.NodeRenderer, len(g.wrappedNodes)+len(g.bindingNodes))
	copy(nodes, g.wrappedNodes)

	for i, n := range g.bindingNodes {
		nodes[len(g.wrappedNodes)+i] = n
	}

	return nodes
}
