package constants

import (
	"fmt"

	"github.com/alkemics/goflow"
	"golang.org/x/tools/go/packages"
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

func craftPkgError(pkgPath string, errs []packages.Error) PkgError {
	es := make([]error, len(errs))
	for i, err := range errs {
		es[i] = err
	}

	return PkgError{
		PkgPath: pkgPath,
		Err:     goflow.MultiError{Errs: es},
	}
}

type TypeError struct {
	Name string
	Err  error
}

func (e TypeError) Error() string {
	return fmt.Sprintf("could not find type of constant %s: %v", e.Name, e.Err)
}

func (e TypeError) Unwrap() error {
	return e.Err
}
