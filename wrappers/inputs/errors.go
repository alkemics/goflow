package inputs

import "fmt"

type BindingError struct {
	Input string
	Err   error
}

func (e BindingError) Error() string {
	return fmt.Sprintf("invalid input %s: %v (inputs)", e.Input, e.Err)
}
