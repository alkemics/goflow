// This package is made of utiliy functions to use goflow with Go.
package gfgo

import (
	"io"

	"github.com/alkemics/goflow"
)

// A GenerateFunc writes g via w. This package provides the utilities to
// generate a graph in Golang.
type GenerateFunc func(w io.Writer, g goflow.GraphRenderer) error
