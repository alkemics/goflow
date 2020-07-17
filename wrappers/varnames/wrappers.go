package varnames

import (
	"github.com/alkemics/goflow"
	"github.com/alkemics/lib-go/v9/sets"
)

func Wrapper(unmarshal func(interface{}) error, graph goflow.GraphRenderer) (goflow.GraphRenderer, error) {
	nodes := graph.Nodes()
	nodeIDs := sets.NewStrings()
	for _, n := range nodes {
		nodeIDs.Add(n.ID())
	}
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
	nodeIDs := sets.NewStrings()
	for _, n := range nodes {
		nodeIDs.Add(n.ID())
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
