package ctx

import (
	"fmt"

	"github.com/alkemics/goflow"
)

const ctxNodeID = "__ctx"

type graphRenderer struct {
	goflow.GraphRenderer

	imports []goflow.Import
	inputs  []goflow.Field
	nodes   []goflow.NodeRenderer
}

func (g graphRenderer) Imports() []goflow.Import     { return g.imports }
func (g graphRenderer) Inputs() []goflow.Field       { return g.inputs }
func (g graphRenderer) Nodes() []goflow.NodeRenderer { return g.nodes }

type nodeRenderer struct {
	goflow.NodeRenderer

	inputs []goflow.Field
}

func (g nodeRenderer) Inputs() []goflow.Field { return g.inputs }

type ctxNode struct{}

func (n ctxNode) ID() string                   { return ctxNodeID }
func (n ctxNode) Previous() []string           { return nil }
func (n ctxNode) Imports() []goflow.Import     { return []goflow.Import{{Pkg: "context", Dir: "context"}} }
func (n ctxNode) Doc() string                  { return "" }
func (n ctxNode) Dependencies() []goflow.Field { return nil }
func (n ctxNode) Inputs() []goflow.Field       { return nil }
func (n ctxNode) Outputs() []goflow.Field {
	return []goflow.Field{{Name: "ctx", Type: "context.Context"}}
}

func (n ctxNode) Run(_, outputs []goflow.Field) (string, error) {
	// ctx must come from the inputs
	return fmt.Sprintf("%s = ctx", outputs[0].Name), nil
}

func Wrapper(_ func(interface{}) error, graph goflow.GraphRenderer) (goflow.GraphRenderer, error) {
	nodes := graph.Nodes()
	wrappedNodes := make([]goflow.NodeRenderer, len(nodes))
	addCtx := false
	for i, node := range nodes {
		inputs := node.Inputs()
		if len(inputs) == 0 || inputs[0].Name != "ctx" || inputs[0].Type != "context.Context" {
			wrappedNodes[i] = node
			continue
		}

		// Here, we will need to add the context in there. We also need to remap the
		// ctx input to inputs.ctx
		addCtx = true
		nodeInputs := make([]goflow.Field, len(inputs))
		copy(nodeInputs, inputs)
		nodeInputs[0].Name = fmt.Sprintf("%s.ctx", ctxNodeID)
		wrappedNodes[i] = nodeRenderer{
			NodeRenderer: node,
			inputs:       nodeInputs,
		}
	}

	inputs := graph.Inputs()
	imports := graph.Imports()
	if addCtx && (len(inputs) == 0 || inputs[0].Name != "ctx" || inputs[0].Type != "context.Context") {
		inputs = append(
			[]goflow.Field{{Name: "ctx", Type: "context.Context"}},
			inputs...,
		)

		imports = append(imports, goflow.Import{Pkg: "context", Dir: "context"})
		wrappedNodes = append(wrappedNodes, ctxNode{})
	}

	return graphRenderer{
		GraphRenderer: graph,

		imports: imports,
		inputs:  inputs,
		nodes:   wrappedNodes,
	}, nil
}
