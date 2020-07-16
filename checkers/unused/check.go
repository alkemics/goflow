package unused

import (
	"fmt"
	"strings"

	"github.com/alkemics/goflow"
)

type Error struct {
	NID string
}

func (e Error) Error() string {
	return fmt.Sprintf("unused node: %s (unused)", e.NID)
}

// Check checks that all nodes are used in one way or another:
// 	- have only one error output
//	- are used explicitly
func Check(graph goflow.GraphRenderer) error {
	nodes := graph.Nodes()

	nextMap := make(map[string][]string)
	for _, n := range nodes {
		for _, p := range n.Previous() {
			nextMap[p] = append(nextMap[p], n.ID())
		}
	}

	errs := make([]error, 0, len(nodes))
	for _, n := range nodes {
		outputs := n.Outputs()
		if strings.HasPrefix(n.ID(), "__") {
			// Internal node.
			continue
		}
		if len(nextMap[n.ID()]) > 0 {
			// Node is explicitly used.
			continue
		}
		if len(outputs) == 0 || len(outputs) == 1 && outputs[0].Type == "error" {
			// Only error output, ignore this node.
			continue
		}
		errs = append(errs, Error{n.ID()})
	}

	if len(errs) > 0 {
		return goflow.MultiError{Errs: errs}
	}

	return nil
}
