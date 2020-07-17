package outputs

import (
	"fmt"
	"strings"
)

type TooManyErrorOutputsError struct {
	Names []string
}

func (e TooManyErrorOutputsError) Error() string {
	return fmt.Sprintf("more than one error output [%s]", strings.Join(e.Names, ", "))
}
