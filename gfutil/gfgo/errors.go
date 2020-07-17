package gfgo

import (
	"fmt"
	"go/types"
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

type TypeError struct {
	Type types.Type
}

func (e TypeError) Error() string {
	return fmt.Sprintf("could no find type of %T", e.Type)
}

type InputParsingError struct {
	InputIndex int
	Err        error
}

func (e InputParsingError) Error() string {
	return fmt.Sprintf("could not read type of input #%d: %v", e.InputIndex, e.Err)
}

func (e InputParsingError) Unwrap() error {
	return e.Err
}
