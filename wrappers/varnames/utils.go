package varnames

import (
	"strings"
)

func compilableGenerateVariableName(s string, nodeIDs map[string]struct{}) string {
	if !strings.Contains(s, ".") {
		return s
	}
	ss := strings.SplitN(s, ".", 2)
	if _, ok := nodeIDs[ss[0]]; !ok {
		return s
	}
	return strings.Join(ss, "_")
}
