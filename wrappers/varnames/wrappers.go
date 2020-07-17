package varnames

import (
	"github.com/alkemics/goflow"
)

func Wrapper(unmarshal func(interface{}) error, graph goflow.GraphRenderer) (goflow.GraphRenderer, error) {
	nodes := graph.Nodes()

	wrappedNodes := make([]goflow.NodeRenderer, len(nodes))
	for i, n := range nodes {
		wrappedNodes[i] = nodeRenderer{NodeRenderer: n}
	}

	return graphRenderer{
		GraphRenderer: graph,
		nodes:         wrappedNodes,
	}, nil
}

func CompilableWrapper(unmarshal func(interface{}) error, graph goflow.GraphRenderer) (goflow.GraphRenderer, error) {
	nodes := graph.Nodes()
	nodeIDs := make(map[string]struct{})
	for _, n := range nodes {
		nodeIDs[n.ID()] = struct{}{}
	}

	wrappedNodes := make([]goflow.NodeRenderer, len(nodes))
	for i, n := range nodes {
		wrappedNodes[i] = compilableNodeRenderer{NodeRenderer: n, nodeIDs: nodeIDs}
	}

	return compilableGraphRenderer{
		GraphRenderer: graph,
		nodes:         wrappedNodes,
		nodeIDs:       nodeIDs,
	}, nil
}
