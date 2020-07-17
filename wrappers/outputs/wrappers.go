package outputs

import (
	"fmt"
	"sort"
	"strings"

	"github.com/alkemics/goflow"
)

func Wrapper(unmarshal func(interface{}) error, graph goflow.GraphRenderer) (goflow.GraphRenderer, error) {
	var outputs struct {
		Outputs map[string]bindings
	}
	if err := unmarshal(&outputs); err != nil {
		return nil, err
	}

	outputBuilders := make([]goflow.NodeRenderer, 0)
	outputFields := make([]goflow.Field, 0, len(outputs.Outputs))
	for output, bs := range outputs.Outputs {
		fields := make([]goflow.Field, len(bs))
		for i, binding := range bs {
			fields[i] = goflow.Field{
				Name: binding,
				Type: "@any",
			}
		}

		outputBuilderName := fmt.Sprintf("__output_%s_builder", output)

		outputFields = append(outputFields, goflow.Field{
			Name: fmt.Sprintf("%s.%s", outputBuilderName, output),
			Type: "@any",
		})

		outputType := fmt.Sprintf("@from[%s]", strings.Join(goflow.FieldNames(fields), ","))

		outputBuilders = append(outputBuilders, outputNode{
			id: fmt.Sprintf("__output_%s_builder", output),
			output: goflow.Field{
				Name: output,
				Type: outputType,
			},
			inputs: fields,
		})
	}

	return graphRenderer{
		GraphRenderer: graph,

		outputs:        outputFields,
		outputBuilders: outputBuilders,
	}, nil
}

func NameWrapper(_ func(interface{}) error, graph goflow.GraphRenderer) (goflow.GraphRenderer, error) {
	outputs := graph.Outputs()
	renamedOutputs := make([]goflow.Field, len(outputs))
	for i, output := range outputs {
		if matches := outputBuilderRegex.FindStringSubmatch(output.Name); len(matches) > 0 {
			output.Name = matches[1]
		}
		renamedOutputs[i] = output
	}

	sort.SliceStable(renamedOutputs, func(i, j int) bool {
		return strings.Compare(renamedOutputs[i].Name, renamedOutputs[j].Name) <= 0
	})

	// Check if at most one error field is returned.
	errorFieldNames := make([]string, 0, len(outputs))
	for _, o := range outputs {
		if o.Type == "error" {
			errorFieldNames = append(errorFieldNames, o.Name)
		}
	}
	if len(errorFieldNames) > 1 {
		return nil, TooManyErrorOutputsError{
			Names: errorFieldNames,
		}
	}

	return nameGraphWrapper{
		GraphRenderer: graph,
		outputs:       renamedOutputs,
	}, nil
}
