package types

import (
	"fmt"
	"strings"

	"github.com/alkemics/goflow"
)

func Wrapper(_ func(interface{}) error, graph goflow.GraphRenderer) (goflow.GraphRenderer, error) {
	// By first sorting the nodes by execution order, we make sure that
	// the output types used to compute an input type have already been
	// resolved and are therefore set properly.
	wrappedNodes := sortNodesByExecutionOrder(graph.Nodes())

	possibleTypes := make(map[string][]string)
	varType := make(map[string]string)

	// Start by storing the outputs: they hold the truth
	for _, n := range wrappedNodes {
		for _, field := range n.Outputs() {
			if field.Type != "" && !strings.HasPrefix(field.Type, "?") && !strings.HasPrefix(field.Type, "@") {
				varType[field.Name] = field.Type
			} else {
				possibleTypes[field.Name] = append(possibleTypes[field.Name], field.Type)
			}
		}
	}

	for _, n := range wrappedNodes {
		for _, field := range n.Inputs() {
			if typ, ok := varType[field.Name]; ok {
				possibleTypes[field.Name] = append(possibleTypes[field.Name], typ)
			} else {
				possibleTypes[field.Name] = append(possibleTypes[field.Name], field.Type)
			}
		}
	}

	// Continue running while we reduce the number of errors.
	iter := 0
	changed := true
	for changed {
		iter++
		if iter > 100 {
			return nil, fmt.Errorf(
				"reached type resolution iteration %d, there probably is an issue with the resolver",
				iter,
			)
		}

		possibleTypes = reduceTypes(possibleTypes)

		before := len(varType)
		for v, typs := range possibleTypes {
			if _, ok := varType[v]; ok {
				continue
			}

			combined := combineTypes(typs)
			if len(combined) == 1 && isTypeResolved(combined[0]) {
				varType[v] = combined[0]
			}
		}

		changed = before < len(varType)
	}

	missings := make(map[string][]string)
	for k, v := range possibleTypes {
		if _, ok := varType[k]; !ok {
			missings[k] = v
		}
	}
	if len(missings) > 0 {
		return graph, goflow.MultiError{
			Errs: craftResolutionErrors(missings),
		}
	}

	for k, n := range wrappedNodes {
		inputs := n.Inputs()
		for i := range inputs {
			field := inputs[i]
			if typ, ok := varType[field.Name]; ok {
				field.Type = typ
			}
			inputs[i] = field
		}

		outputs := n.Outputs()
		for i := range outputs {
			field := outputs[i]
			if typ, ok := varType[field.Name]; ok {
				field.Type = typ
			}
			outputs[i] = field
		}

		wrappedNodes[k] = nodeRenderer{
			NodeRenderer: n,
			typedInputs:  inputs,
			typedOutputs: outputs,
		}
	}

	return graphRenderer{
		GraphRenderer: graph,

		nodes:     wrappedNodes,
		outputMap: varType,
	}, nil
}
