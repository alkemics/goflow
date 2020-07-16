package constants

import (
	"fmt"

	"github.com/alkemics/goflow"
)

func nodeWrapper(constants []constant) goflow.NodeWrapper {
	return func(_ func(interface{}) error, node goflow.NodeRenderer) (goflow.NodeRenderer, error) {
		inputs := node.Inputs()
		parsedInputs := make([]goflow.Field, len(inputs))
		copy(parsedInputs, inputs)

		extraImports := make([]goflow.Import, 0)
		for i, input := range inputs {
			cst := findConstant(input.Name, constants)
			if cst.name == "" {
				cst = findHardcodedValue(input.Name)
			}

			if cst.name == "" {
				// It's not a constant.
				continue
			}

			typ := cst.typ
			parsedInputs[i] = goflow.Field{
				Name: input.Name,
				Type: fmt.Sprintf("@type[%s,%s]", typ, input.Type),
			}

			if cst.imp.Dir != "" {
				extraImports = append(extraImports, cst.imp)
			}
		}

		return nodeRenderer{
			NodeRenderer: node,
			imports:      extraImports,
			inputs:       parsedInputs,
		}, nil
	}
}

func Wrapper(constantPackages ...string) goflow.GraphWrapper {
	constants, err := loadConstants(constantPackages)
	return func(unmarshal func(interface{}) error, graph goflow.GraphRenderer) (goflow.GraphRenderer, error) {
		if err != nil {
			return nil, goflow.GraphError{
				Wrapper: "constants",
				Err:     fmt.Errorf("could not load constants: %w", err),
			}
		}

		wrapper := goflow.FromNodeWrapper(nodeWrapper(constants))
		return wrapper(unmarshal, graph)
	}
}
