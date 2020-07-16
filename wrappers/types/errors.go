package types

import (
	"fmt"
	"strings"
)

type ResolutionError struct {
	NID   string
	Var   string
	Types []string
}

func (e ResolutionError) Error() string {
	return fmt.Sprintf("%s.%s: could not resolve type from [%s] (types)", e.NID, e.Var, strings.Join(e.Types, ", "))
}

func (e ResolutionError) NodeID() string  { return e.NID }
func (e ResolutionError) Wrapper() string { return "types" }

func craftResolutionErrors(types map[string][]string) []error {
	errs := make([]error, 0, len(types))
	for key, ts := range types {
		nodeID := key
		v := "<unkown>"
		if strings.Contains(key, ".") {
			split := strings.SplitN(key, ".", 2)
			nodeID, v = split[0], split[1]
		}
		errs = append(errs, ResolutionError{
			NID:   nodeID,
			Var:   v,
			Types: ts,
		})
	}
	return errs
}
