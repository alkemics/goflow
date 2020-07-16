package bind

import (
	"fmt"
	"strings"

	"github.com/alkemics/goflow"
)

func Wrapper(unmarshal func(interface{}) error, graph goflow.GraphRenderer) (goflow.GraphRenderer, error) {
	var binds struct {
		Nodes []struct {
			ID   string
			Bind map[string]bindings
		}
	}
	if err := unmarshal(&binds); err != nil {
		return nil, err
	}

	bindingsByID := make(map[string]map[string]bindings)
	for _, node := range binds.Nodes {
		bindingsByID[node.ID] = node.Bind
	}

	errs := make([]error, 0)
	nodes := graph.Nodes()
	wrappedNodes := make([]goflow.NodeRenderer, len(nodes))
	bindingNodes := make([]bindingNode, 0)
	for i, node := range nodes {
		bn := bindingsByID[node.ID()]
		if bn == nil {
			bn = make(map[string]bindings)
		}
		nodeInputs := node.Inputs()

		unkownInputs := make([]string, 0)
		bindMap := make(map[string]goflow.Field)
		for input, bs := range bn {
			var nodeInput goflow.Field
			for _, in := range nodeInputs {
				if in.Name == input {
					nodeInput = in
					break
				}
			}

			if nodeInput.Name == "" {
				unkownInputs = append(unkownInputs, input)
				continue
			}

			aggregatorNodeID := fmt.Sprintf("__%s_%s", node.ID(), input)
			inputs := make([]goflow.Field, len(bs))
			for i, b := range bs {
				f := goflow.Field{
					Name: b,
					Type: fmt.Sprintf("?%s", nodeInput.Type),
				}
				inputs[i] = f
			}
			bindingNodes = append(bindingNodes, bindingNode{
				id:     aggregatorNodeID,
				inputs: inputs,
				output: goflow.Field{
					Name: "aggregated",
					Type: nodeInput.Type,
				},
			})
			bindMap[input] = goflow.Field{
				Name: fmt.Sprintf("%s.aggregated", aggregatorNodeID),
				Type: nodeInput.Type,
			}
		}

		if len(unkownInputs) > 0 {
			errs = append(errs, goflow.NodeError{
				ID:      node.ID(),
				Wrapper: "bind",
				Err:     fmt.Errorf("unknown inputs: [%s]", strings.Join(unkownInputs, ", ")),
			})
		}

		// Check that all inputs have been given a value.
		missingInputNames := make([]string, 0)
		optionalInputNames := getOptionalInputNames(node.Doc())
		absentInputNames := make(map[string]struct{})
		for _, input := range nodeInputs {
			// Needed to allow other wrappers to inject inputs
			// TODO: inject the decorators via inputs.decorators
			if strings.Contains(input.Name, ".") || input.Name == "decorators" {
				continue
			}

			if _, ok := bindMap[input.Name]; !ok {
				if _, ok := optionalInputNames[input.Name]; ok {
					absentInputNames[input.Name] = struct{}{}
				} else {
					missingInputNames = append(missingInputNames, input.Name)
				}
			}
		}
		if len(missingInputNames) > 0 {
			errs = append(errs, goflow.NodeError{
				ID:      node.ID(),
				Wrapper: "bind",
				Err:     fmt.Errorf("inputs not bound: [%s]", strings.Join(missingInputNames, ", ")),
			})
			continue
		}

		wrappedNodes[i] = nodeRenderer{
			NodeRenderer:     node,
			bindMap:          bindMap,
			absentInputNames: absentInputNames,
		}
	}

	if len(errs) > 0 {
		return nil, goflow.MultiError{Errs: errs}
	}

	return graphRenderer{
		GraphRenderer: graph,
		wrappedNodes:  wrappedNodes,
		bindingNodes:  bindingNodes,
	}, nil
}
