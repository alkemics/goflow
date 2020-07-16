package gfgo

import (
	"fmt"
)

type PkgError struct {
	PkgPath string
	Err     error
}

func (e PkgError) Error() string {
	return fmt.Sprintf("loading %s: %v", e.PkgPath, e.Err)
}

func (e PkgError) Unwrap() error {
	return e.Err
}
