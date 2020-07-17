package bind

import (
	"fmt"
	"strings"

	"github.com/alkemics/goflow"
)

type bindings []string

func (b *bindings) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multipleBindings []string
	if err := unmarshal(&multipleBindings); err == nil {
		*b = multipleBindings
		return nil
	}

	var simpleBinding string
	if err := unmarshal(&simpleBinding); err != nil {
		return err
	}

	*b = []string{simpleBinding}
	return nil
}

func getOptionalInputNames(doc string) map[string]struct{} {
	optionalInputNames := make(map[string]struct{})
	for _, line := range strings.Split(doc, "\n") {
		if !strings.HasPrefix(line, "optional:") {
			continue
		}
		for _, name := range strings.Split(strings.TrimPrefix(line, "optional:"), ",") {
			optionalInputNames[strings.Trim(name, " ")] = struct{}{}
		}
	}
	return optionalInputNames
}

func Run(inputs, outputs []goflow.Field) (string, error) {
	outputName := outputs[0].Name
	outputType := outputs[0].Type
	runs := make([]string, len(inputs))

	if len(inputs) == 1 {
		input := inputs[0]
		if input.Type == outputType {
			return fmt.Sprintf("%s = %s", outputName, input.Name), nil
		}

		if strings.TrimPrefix(outputType, "[]") == input.Type {
			// outputType is [][]T and input.Type is []T: append is the way to go
			return fmt.Sprintf("%s = append(%s, %s)", outputName, outputName, input.Name), nil
		}

		if strings.HasPrefix(input.Type, "[]") {
			// Both are slices but not exactly the same
			return fmt.Sprintf(`for _, e := range %s {
				%s = append(%s, e)
			}`, input.Name, outputName, outputName), nil
		}

		if strings.HasPrefix(outputType, "[]") {
			return fmt.Sprintf("%s = append(%s, %s)", outputName, outputName, input.Name), nil
		}
		return fmt.Sprintf("%s = %s", outputName, input.Name), nil
	}

	for i, input := range inputs {
		if strings.HasPrefix(outputType, "[]") {
			if input.Type == outputType {
				// both are slices of the same type: append is the way to go
				runs[i] = fmt.Sprintf("%s = append(%s, %s...)", outputName, outputName, input.Name)
				continue
			}

			if strings.TrimPrefix(outputType, "[]") == input.Type {
				// outputType is [][]T and input.Type is []T: append is the way to go
				runs[i] = fmt.Sprintf("%s = append(%s, %s)", outputName, outputName, input.Name)
				continue
			}

			if strings.HasPrefix(input.Type, "[]") {
				// Both are slices but not exactly the same
				runs[i] = fmt.Sprintf(`for _, e := range %s {
					%s = append(%s, e)
				}`, input.Name, outputName, outputName)
				continue
			}

			runs[i] = fmt.Sprintf("%s = append(%s, %s)", outputName, outputName, input.Name)
			continue
		}

		runs[i] = fmt.Sprintf("%s = %s", outputName, input.Name)
	}

	return strings.Join(runs, "\n"), nil
}
