package previous

import (
	"fmt"

	"github.com/alkemics/goflow"
)

type Error struct {
	NID string
}

func (e Error) Error() string {
	return fmt.Sprintf("node not found: %s (after)", e.NID)
}

// Check checks that all nodes in previous exist.
func Check(graph goflow.GraphRenderer) error {
	nodes := graph.Nodes()

	allNodeIDs := make(map[string]struct{})
	for _, n := range nodes {
		allNodeIDs[n.ID()] = struct{}{}
	}

	errs := make([]error, 0)
	for _, n := range nodes {
		for _, p := range n.Previous() {
			if _, ok := allNodeIDs[p]; !ok {
				errs = append(errs, Error{p})
			}
		}
	}

	if len(errs) > 0 {
		return goflow.MultiError{Errs: errs}
	}

	return nil
}
