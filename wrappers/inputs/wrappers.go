package inputs

import (
	"strings"

	"github.com/alkemics/goflow"
)

func Wrapper(unmarshal func(interface{}) error, graph goflow.GraphRenderer) (goflow.GraphRenderer, error) {
	var inputs struct {
		Inputs []binding `yaml:"inputs"`
	}
	if err := unmarshal(&inputs); err != nil {
		return nil, err
	}

	fields := make([]goflow.Field, len(inputs.Inputs))
	optionalInputNames := make([]string, 0, len(inputs.Inputs))
	for i, f := range inputs.Inputs {
		if strings.HasSuffix(f.Name, "?") {
			f = binding(goflow.Field{
				Name: strings.TrimSuffix(f.Name, "?"),
				Type: f.Type,
			})
			optionalInputNames = append(optionalInputNames, f.Name)
		}
		fields[i] = goflow.Field(f)
	}

	return graphRenderer{
		GraphRenderer: graph,

		inputs:             fields,
		inputNode:          inputNode{outputs: fields},
		optionalInputNames: optionalInputNames,
	}, nil
}

func TypeWrapper(_ func(interface{}) error, graph goflow.GraphRenderer) (goflow.GraphRenderer, error) {
	return graphTypeRenderer{
		GraphRenderer: graph,
	}, nil
}
