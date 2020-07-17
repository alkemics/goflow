package gonodes

import "fmt"

type NotFoundError struct {
	ID   string
	Type string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s: unkown node type: %s (gonodes)", e.ID, e.Type)
}
