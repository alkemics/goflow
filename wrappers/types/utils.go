package types

import (
	"fmt"
	"sort"
	"strings"

	"github.com/alkemics/goflow"
	"github.com/alkemics/goflow/gfutil"
	"github.com/alkemics/lib-go/v9/sets"
)

var numberTypes = sets.NewStrings(
	"int",
	"int8",
	"int16",
	"int32",
	"int64",
	"uint",
	"uint8",
	"uint16",
	"uint32",
	"uint64",
	"float32",
	"float64",
)

func isEmptyInterface(s string) bool {
	return s == "interface{}" || s == "[]interface{}"
}

func isSlice(s string) bool {
	return strings.HasPrefix(s, "[]")
}

func trimSlice(s string) string {
	return strings.TrimPrefix(s, "[]")
}

func makeSlice(s string) string {
	return fmt.Sprintf("[]%s", s)
}

func ensureSlice(s string) string {
	if isSlice(s) {
		return s
	}
	return makeSlice(s)
}

func trimOption(s string) string {
	return strings.TrimPrefix(s, "?")
}

func isTypeResolved(typ string) bool {
	return typ != "" && !strings.HasPrefix(typ, "@") && !strings.HasPrefix(typ, "?")
}

func fixType(t string) string {
	if strings.Contains(t, "@any") {
		return "@any"
	}

	if strings.HasPrefix(t, "[]?") {
		return fmt.Sprintf("?[]%s", t[3:])
	}

	return t
}

func reduceTypes(possibleTypes map[string][]string) map[string][]string {
	agg := make(map[string][]string)
	for n, types := range possibleTypes {
		for _, typ := range types {
			if strings.HasPrefix(typ, "@from[") && strings.HasSuffix(typ, "]") {
				names := strings.Split(strings.TrimPrefix(strings.TrimSuffix(typ, "]"), "@from["), ",")
				needSlice := len(names) > 1
				for _, name := range names {
					ts := make([]string, len(possibleTypes[name]))
					copy(ts, possibleTypes[name])
					for i, t := range ts {
						if needSlice {
							t = ensureSlice(t)
						}
						ts[i] = fixType(t)
					}
					agg[n] = append(agg[n], ts...)
				}
			} else if strings.HasPrefix(typ, "@type[") && strings.HasSuffix(typ, "]") {
				ts := strings.Split(strings.TrimPrefix(strings.TrimSuffix(typ, "]"), "@type["), ",")
				agg[n] = append(agg[n], ts...)
			} else if typ == "[]@single" {
				agg[n] = append(agg[n], "@any")
			} else if typ == "untyped int" {
				agg[n] = append(agg[n], "@number")
			} else if typ == "untyped float" {
				agg[n] = append(agg[n], "@number")
			} else if typ == "untyped string" {
				agg[n] = append(agg[n], "?string")
			} else {
				agg[n] = append(agg[n], typ)
			}
		}

		agg[n] = combineTypes(agg[n])
	}
	return agg
}

func combineTypes(types []string) []string {
	if len(types) == 0 {
		return nil
	}

	combined := sets.NewStrings(types...)

	// If we have @any and any other element(s), we keep the other element(s)
	if combined.Contains("@any") && combined.Len() > 1 {
		combined.Remove("@any")
	}

	flat := combined.ToList()
	current := flat[0]
	combined = sets.NewStrings()
	for _, next := range flat[1:] {
		if current == "" {
			current = next
			continue
		}

		// We do not need to check for current == next because
		// flat comes from a sets.Strings!

		cc := trimOption(current)
		nn := trimOption(next)
		if cc == nn {
			// The types are not EXACTLY equal, but with no option
			// they are: this means that 'current' is ?T and 'next' is T
			// or vice-versa -> we keep T
			current = cc
			continue
		}
		// Here the types are not equal even with the option trimmed

		if isEmptyInterface(cc) && !isEmptyInterface(nn) {
			current = nn
			continue
		}
		if !isEmptyInterface(cc) && isEmptyInterface(nn) {
			current = cc
			continue
		}
		if isEmptyInterface(cc) && isEmptyInterface(nn) {
			continue
		}
		// We are now certain none of the types are empty interfaces

		if trimSlice(cc) == trimSlice(nn) {
			// The underlying type of cc and nn are equal -> we keep it
			current = trimSlice(cc)
			continue
		}
		if trimSlice(cc) == nn {
			current = nn
			continue
		}
		if trimSlice(nn) == cc {
			current = cc
			continue
		}
		// And now, we have no more slices

		// Check the additional directives
		// @number directive: if the other type is a number, all is good!
		if cc == "@number" && numberTypes.Contains(trimSlice(nn)) {
			current = trimSlice(nn)
			continue
		}
		if nn == "@number" && numberTypes.Contains(trimSlice(cc)) {
			current = trimSlice(cc)
			continue
		}

		if cc == "@single" {
			current = nn
			if isSlice(nn) {
				current = trimSlice(nn)
			}
			continue
		}
		if nn == "@single" {
			current = cc
			if isSlice(cc) {
				current = trimSlice(cc)
			}
			continue
		}

		// ?bool is basically with everything
		if current == "?bool" {
			current = next
			continue
		}
		if next == "?bool" {
			continue
		}

		// We could not match the types so we add the current one
		// as we need more resolution/information to match the types
		combined.Add(current)
		current = next
	}
	combined.Add(current)

	if combined.Len() == 1 {
		// We have only one element, let's resolve it!
		typ := combined.ToList()[0]
		return []string{trimOption(typ)}
	}

	return combined.ToList()
}

func sortNodesByExecutionOrder(nodes []goflow.NodeRenderer) []goflow.NodeRenderer {
	sortedNodes := make([]goflow.NodeRenderer, len(nodes))
	copy(sortedNodes, nodes)

	dependencies := gfutil.UnravelNodeDependencies(sortedNodes)
	sort.SliceStable(sortedNodes, func(i, j int) bool {
		firstNodeID := sortedNodes[i].ID()
		secondNodeID := sortedNodes[j].ID()

		for _, dep := range dependencies[firstNodeID] {
			if dep == secondNodeID {
				return false
			}
		}
		return true
	})

	return sortedNodes
}
