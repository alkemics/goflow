package cycles

import (
	"fmt"

	"github.com/alkemics/goflow"
)

type CycleError struct {
	NID string
}

func (e CycleError) Error() string {
	return fmt.Sprintf("node visited multiple times: %s (cycles)", e.NID)
}

type OrphanError struct {
	NID string
}

func (e OrphanError) Error() string {
	return fmt.Sprintf("orphan node: %s (cycles)", e.NID)
}

// Check checks that no cycles exist in the graph as well as for
// orphan nodes, i.e. nodes that are never visited
func Check(graph goflow.GraphRenderer) error {
	// Map nodes by ID.
	nodeMap := make(map[string]goflow.NodeRenderer)
	for _, n := range graph.Nodes() {
		nodeMap[n.ID()] = n
	}

	// Construct a map of edges and find the nodes without parents.
	firstNodeIDs := make([]string, 0)
	edges := make(map[string][]string)
	for nID, n := range nodeMap {
		previous := n.Previous()
		if len(previous) == 0 {
			firstNodeIDs = append(firstNodeIDs, nID)
		}
		for _, prev := range previous {
			edges[prev] = append(edges[prev], nID)
		}
	}

	// v, ok:
	// ok == false -> not visited
	// v == true -> permanent mark
	// v == false -> temporary mark
	visited := make(map[string]bool)

	for _, nID := range firstNodeIDs {
		if err := visit(nID, edges, visited); err != nil {
			return err
		}
	}

	// Check that all nodes have been visited.
	orphanErrs := make([]error, 0)
	for nID := range nodeMap {
		permanent := visited[nID]
		if !permanent {
			orphanErrs = append(orphanErrs, OrphanError{NID: nID})
		}
	}
	if len(orphanErrs) > 0 {
		return goflow.MultiError{Errs: orphanErrs}
	}

	return nil
}

// visit implements the DFS algorithm presented here
// https://en.wikipedia.org/wiki/Topological_sorting#Depth-first_search
func visit(nodeID string, edges map[string][]string, visited map[string]bool) error {
	if visited == nil {
		return nil
	}
	permanent, marked := visited[nodeID]
	if marked {
		if permanent {
			return nil
		}
		return CycleError{NID: nodeID}
	}

	visited[nodeID] = false

	for _, nextID := range edges[nodeID] {
		if err := visit(nextID, edges, visited); err != nil {
			return err
		}
	}

	visited[nodeID] = true
	return nil
}
