package gfutil

import (
	"sort"

	"github.com/alkemics/goflow"
)

func UnravelDependencies(dependencies map[string][]string) map[string][]string {
	deps := make(map[string]map[string]struct{})
	for k, v := range dependencies {
		deps[k] = make(map[string]struct{})
		for _, s := range v {
			deps[k][s] = struct{}{}
		}
	}

	var (
		allDependencies func(k string, dependencies map[string]map[string]struct{}) map[string]struct{}
		handledNodes    = make(map[string]struct{})
	)

	allDependencies = func(k string, dependencies map[string]map[string]struct{}) map[string]struct{} {
		allDeps := dependencies[k]

		if _, ok := handledNodes[k]; !ok {
			handledNodes[k] = struct{}{}

			for dep := range allDeps {
				for other := range allDependencies(dep, dependencies) {
					allDeps[other] = struct{}{}
				}
			}

			dependencies[k] = allDeps
		}

		return allDeps
	}

	for k := range deps {
		deps[k] = allDependencies(k, deps)
	}

	unraveled := make(map[string][]string)

	for k, v := range deps {
		if _, ok := dependencies[k]; ok {
			// Add a default ONLY if it was in the dependency tree
			// originally
			unraveled[k] = nil
		}

		for s := range v {
			unraveled[k] = append(unraveled[k], s)
		}

		sort.Strings(unraveled[k])
	}

	return unraveled
}

func UnravelNodeDependencies(nodes []goflow.NodeRenderer) map[string][]string {
	deps := make(map[string][]string)
	for _, node := range nodes {
		deps[node.ID()] = node.Previous()
	}

	return UnravelDependencies(deps)
}

func SortNodes(nodes []goflow.NodeRenderer) {
	dependencies := UnravelNodeDependencies(nodes)
	sort.SliceStable(nodes, func(i, j int) bool {
		firstNodeID := nodes[i].ID()
		secondNodeID := nodes[j].ID()

		// The first node depends on the second -> first > second
		for _, dep := range dependencies[firstNodeID] {
			if dep == secondNodeID {
				return false
			}
		}

		// The second node depends on the first -> first < second
		for _, dep := range dependencies[secondNodeID] {
			if dep == firstNodeID {
				return true
			}
		}

		// The first and seconds have a different number of dependencies
		// -> first < seconds <=> len(first.deps) < len(second.deps)
		if len(dependencies[firstNodeID]) != len(dependencies[secondNodeID]) {
			return len(dependencies[firstNodeID]) < len(dependencies[secondNodeID])
		}

		// If the nodes are independent and have the same number of dependencies
		// -> compare the ids
		return firstNodeID <= secondNodeID
	})
}
