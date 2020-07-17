package varnames

import (
	"fmt"
	"strings"

	"github.com/alkemics/goflow"
)

type nodeRenderer struct {
	goflow.NodeRenderer
}

func (n nodeRenderer) Outputs() []goflow.Field {
	outputs := n.NodeRenderer.Outputs()
	namedOutputs := make([]goflow.Field, len(outputs))
	for i, f := range outputs {
		f.Name = fmt.Sprintf("%s.%s", n.ID(), f.Name)
		namedOutputs[i] = f
	}
	return namedOutputs
}

type graphRenderer struct {
	goflow.GraphRenderer

	nodes []goflow.NodeRenderer
}

func (g graphRenderer) Nodes() []goflow.NodeRenderer { return g.nodes }

// CompilableWrapper
// TODO: find a real name...

type compilableNodeRenderer struct {
	goflow.NodeRenderer

	nodeIDs map[string]struct{}
}

func (n compilableNodeRenderer) Previous() []string {
	previous := n.NodeRenderer.Previous()
	previousSet := make(map[string]struct{})
	for _, p := range previous {
		previousSet[p] = struct{}{}
	}

	inputs := n.NodeRenderer.Inputs()
	for _, f := range inputs {
		nodeName := strings.SplitN(f.Name, ".", 2)[0]
		if _, ok := n.nodeIDs[nodeName]; !ok {
			continue
		}

		if _, ok := previousSet[nodeName]; ok {
			// If we already know that node, skip it
			continue
		}

		previousSet[nodeName] = struct{}{}
		previous = append(previous, nodeName)
	}
	return previous
}

func (n compilableNodeRenderer) Inputs() []goflow.Field {
	inputs := n.NodeRenderer.Inputs()
	namedInputs := make([]goflow.Field, len(inputs))
	for i, f := range inputs {
		f.Name = compilableGenerateVariableName(f.Name, n.nodeIDs)
		namedInputs[i] = f
	}
	return namedInputs
}

func (n compilableNodeRenderer) Outputs() []goflow.Field {
	outputs := n.NodeRenderer.Outputs()
	namedOutputs := make([]goflow.Field, len(outputs))
	for i, f := range outputs {
		f.Name = compilableGenerateVariableName(f.Name, n.nodeIDs)
		namedOutputs[i] = f
	}
	return namedOutputs
}

type compilableGraphRenderer struct {
	goflow.GraphRenderer
	nodes   []goflow.NodeRenderer
	nodeIDs map[string]struct{}
}

func (g compilableGraphRenderer) Nodes() []goflow.NodeRenderer { return g.nodes }

func (g compilableGraphRenderer) Outputs() []goflow.Field {
	outputs := g.GraphRenderer.Outputs()
	namedOutputs := make([]goflow.Field, len(outputs))
	for i, f := range outputs {
		f.Name = compilableGenerateVariableName(f.Name, g.nodeIDs)
		namedOutputs[i] = f
	}
	return namedOutputs
}
