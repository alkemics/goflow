package goflow

// The GraphWrapper function method is the main building block of GoFlow. It is
// used to add features to goflow, such as inspection or keywords.
//
// Implementing a GraphWrapper might not always be the best fit. If you need something
// simpler, make sure you have a look at the other blocks:
//   - NodeWrapper
//   - Linter
//   - Checker
// goflow provides functions to use those simpler types as GraphWrappers.
type GraphWrapper func(unmarshal func(interface{}) error, graph GraphRenderer) (GraphRenderer, error)

// The NodeWrapper should be used when no operation needs to be done on
// the graph itself, i.e. when everything is handled at the node level.
//
// To register a NodeWrapper in the builder, use FromNodeWrapper.
type NodeWrapper func(unmarshal func(interface{}) error, node NodeRenderer) (NodeRenderer, error)

type nodeWrapperGraphRenderer struct {
	GraphRenderer

	nodes []NodeRenderer
}

func (g nodeWrapperGraphRenderer) Nodes() []NodeRenderer {
	return g.nodes
}

func (g nodeWrapperGraphRenderer) Imports() []Import {
	imports := g.GraphRenderer.Imports()

	for _, node := range g.Nodes() {
		imports = append(imports, node.Imports()...)
	}

	return imports
}

func fakeUnmarshal(interface{}) error { return nil }

// FromNodeWrapper converts a NodeWrapper into a GraphWrapper.
func FromNodeWrapper(nw NodeWrapper) GraphWrapper {
	return func(unmarshal func(interface{}) error, graph GraphRenderer) (GraphRenderer, error) {
		var nodes struct {
			Nodes []nodeLoader
		}
		if err := unmarshal(&nodes); err != nil {
			return nil, err
		}

		ns := graph.Nodes()
		errs := make([]error, 0, len(ns))
		wrappedNodes := make([]NodeRenderer, len(ns))
		copy(wrappedNodes, ns)

		for i, node := range ns {
			node := node
			nodeUnmarshal := fakeUnmarshal
			for _, n := range nodes.Nodes {
				if n.IDVal == node.ID() {
					nodeUnmarshal = n.Unmarshal
					break
				}
			}

			wrapped, err := nw(nodeUnmarshal, node)
			if err != nil {
				errs = append(errs, err)
			}

			if wrapped != nil {
				wrappedNodes[i] = wrapped
			}
		}

		var err error
		if len(errs) > 0 {
			err = MultiError{Errs: errs}
		}

		return nodeWrapperGraphRenderer{
			GraphRenderer: graph,
			nodes:         wrappedNodes,
		}, err
	}
}

// A Linter is used to lint the yaml file to find YAML-related errors, such
// as duplicated ids.
//
// A Linter has access to the unmarshal function so it can focus on the yaml.
// To check errors in the graph, use a Checker.
//
// To use lintes with goflow, the easiest way is to wrap them with FromLinter
// and place them first in your list of wrappers.
type Linter func(unmarshal func(interface{}) error) error

// FromLinter transforms a linter into a GraphWrapper
func FromLinter(lint Linter) GraphWrapper {
	return func(unmarshal func(interface{}) error, graph GraphRenderer) (GraphRenderer, error) {
		return nil, lint(unmarshal)
	}
}

// A Checker is used to check that a graph follows the guide lines you set, such as
// checking that no node is never used.
//
// A Checker has access to the GraphRenderer, making it a good fit for coherence checks
// on the graph. As it does not have access to the unmarshal function, it cannot
// lint the yaml. To do that, use a Linter.
//
// To use checkers with goflow, the easiest way is to wrap them with FromChecker
// and place the last in you list of wrappers.
type Checker func(graph GraphRenderer) error

// FromChecker transforms a checker into a GraphWrapper
func FromChecker(check Checker) GraphWrapper {
	return func(unmarshal func(interface{}) error, graph GraphRenderer) (GraphRenderer, error) {
		if err := check(graph); err != nil {
			return nil, err
		}
		return graph, nil
	}
}
