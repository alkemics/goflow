package ifs

import (
	"fmt"
	"strings"

	"github.com/alkemics/goflow"
)

type condition struct {
	name    string
	negated bool
}

func getConditions(list []string) []condition {
	conditions := make([]condition, len(list))
	for i, e := range list {
		name := e
		negated := false
		if strings.HasPrefix(name, "not ") {
			name = strings.TrimPrefix(name, "not ")
			negated = true
		}
		conditions[i] = condition{
			name:    name,
			negated: negated,
		}
	}
	return conditions
}

func generateConditionals(inputs []goflow.Field, conditions []condition) string {
	conditionals := make([]string, len(inputs))
	for i, input := range inputs {
		negated := conditions[i].negated
		if strings.HasPrefix(input.Type, "[]") {
			sign := ">"
			if negated {
				sign = "=="
			}
			conditionals[i] = fmt.Sprintf("len(%s) %s 0", input.Name, sign)
		} else if input.Type == "error" || strings.HasPrefix(input.Type, "*") {
			sign := "!="
			if negated {
				sign = "=="
			}
			conditionals[i] = fmt.Sprintf("%s %s nil", input.Name, sign)
		} else if isNumber(input.Type) {
			sign := "!="
			if negated {
				sign = "=="
			}
			conditionals[i] = fmt.Sprintf("%s %s 0", input.Name, sign)
		} else if input.Type == "string" {
			sign := "!="
			if negated {
				sign = "=="
			}
			conditionals[i] = fmt.Sprintf("%s %s \"\"", input.Name, sign)
		} else {
			prefix := ""
			if negated {
				prefix = "!"
			}
			conditionals[i] = fmt.Sprintf("%s%s", prefix, input.Name)
		}
	}
	return strings.Join(conditionals, " && ")
}

func isNumber(typ string) bool {
	return typ == "int" || typ == "int8" || typ == "int16" || typ == "int32" || typ == "int64" || typ == "uint" || typ == "uint8" || typ == "uint16" || typ == "uint32" || typ == "uint64" || typ == "float32" || typ == "float64"
}
