package gonodes

import (
	"sort"
	"strings"

	"github.com/alkemics/goflow"
	"github.com/alkemics/goflow/gfutil/gfgo"
)

type NodeFinder interface {
	Find(typ string) gfgo.Node
}

func nodeWrapper(nodeFinder NodeFinder, graphPkg string) goflow.NodeWrapper {
	return func(unmarshal func(interface{}) error, node goflow.NodeRenderer) (goflow.NodeRenderer, error) {
		var typed struct {
			Type string
		}

		if err := unmarshal(&typed); err != nil {
			return nil, err
		}
		typ := typed.Type

		if typ == "" {
			return nil, nil
		}

		// Go find the go node based on `typ`.
		gn := nodeFinder.Find(typ)
		if gn.Typ == "" {
			return nil, NotFoundError{
				ID:   node.ID(),
				Type: typ,
			}
		}

		return nodeRenderer{
			NodeRenderer: node,
			graphPkg:     graphPkg,
			node:         gn,
		}, nil
	}
}

func Wrapper(nodeFinder NodeFinder) goflow.GraphWrapper {
	return func(u func(interface{}) error, g goflow.GraphRenderer) (goflow.GraphRenderer, error) {
		wrapper := goflow.FromNodeWrapper(nodeWrapper(nodeFinder, g.Pkg()))
		return wrapper(u, g)
	}
}

func DepWrapper(unmarshal func(interface{}) error, graph goflow.GraphRenderer) (goflow.GraphRenderer, error) {
	deps := make([]goflow.Field, 0)
	for _, n := range graph.Nodes() {
		deps = mergeDependencies(deps, n.Dependencies()...)
	}

	sort.Slice(deps, func(i, j int) bool {
		return strings.Compare(deps[i].Name, deps[j].Name) < 0
	})
	return depsGraphRenderer{
		GraphRenderer: graph,
		deps:          deps,
	}, nil
}

func mergeDependencies(fields []goflow.Field, others ...goflow.Field) []goflow.Field {
	if len(others) == 0 {
		return fields
	}
	fMap := make(map[goflow.Field]struct{})
	for _, f := range fields {
		fMap[f] = struct{}{}
	}
	for _, f := range others {
		if _, ok := fMap[f]; !ok {
			fields = append(fields, f)
			fMap[f] = struct{}{}
		}
	}
	return fields
}
