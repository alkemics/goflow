package varnames

import (
	"strings"

	"github.com/alkemics/lib-go/v9/sets"
)

func compilableGenerateVariableName(s string, nodeIDs sets.Strings) string {
	if !strings.Contains(s, ".") {
		return s
	}
	ss := strings.SplitN(s, ".", 2)
	if !nodeIDs.Contains(ss[0]) {
		return s
	}
	return strings.Join(ss, "_")
}
