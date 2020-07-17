package gonodes

import (
	"fmt"
	"strings"

	"github.com/alkemics/goflow"
	"github.com/alkemics/goflow/gfutil/gfgo"
)

type nodeRenderer struct {
	goflow.NodeRenderer

	graphPkg string

	node gfgo.Node
}

func (n nodeRenderer) Imports() []goflow.Import {
	return append(n.NodeRenderer.Imports(), n.node.Imports...)
}
func (n nodeRenderer) Doc() string                  { return n.node.Doc }
func (n nodeRenderer) Dependencies() []goflow.Field { return n.node.Dependencies }
func (n nodeRenderer) Inputs() []goflow.Field       { return n.node.Inputs }
func (n nodeRenderer) Outputs() []goflow.Field      { return n.node.Outputs }
func (n nodeRenderer) Run(inputs, outputs []goflow.Field) (string, error) {
	statements := make([]string, 0)
	call := fmt.Sprintf("%s.%s", n.node.Pkg, n.node.Typ)

	if n.node.Method != "" {
		if n.node.Constructor != "" {
			depNames := goflow.FieldNames(n.node.Dependencies)
			for i := range depNames {
				depNames[i] = fmt.Sprintf("g.%s", depNames[i])
			}
			statements = append(statements, fmt.Sprintf(
				"node := %s(%s)",
				n.formatNode(n.node.Pkg, n.node.Constructor),
				strings.Join(depNames, ", "),
			))
		} else {
			statements = append(statements, fmt.Sprintf(
				"node := %s{}",
				n.formatNode(n.node.Pkg, n.node.Constructor),
			))
		}

		call = fmt.Sprintf("node.%s", n.node.Method)
	}

	if len(outputs) == 0 {
		statements = append(statements, fmt.Sprintf(
			"%s(%s)",
			call,
			strings.Join(goflow.FieldNames(inputs), ", "),
		))
	} else {
		statements = append(statements, fmt.Sprintf(
			"%s = %s(%s)",
			strings.Join(goflow.FieldNames(outputs), ", "),
			call,
			strings.Join(goflow.FieldNames(inputs), ", "),
		))
	}

	return strings.Join(statements, "\n"), nil
}

func (n nodeRenderer) formatNode(pkg, name string) string {
	if pkg == n.graphPkg {
		return name
	}
	return fmt.Sprintf("%s.%s", pkg, name)
}

type depsGraphRenderer struct {
	goflow.GraphRenderer

	deps []goflow.Field
}

func (g depsGraphRenderer) Dependencies() []goflow.Field { return g.deps }
