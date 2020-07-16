package cycles_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alkemics/goflow"
	"github.com/alkemics/goflow/checkers/cycles"
	"github.com/alkemics/goflow/gfutil/gfgo"
	"github.com/alkemics/goflow/wrappers/after"
	"github.com/alkemics/goflow/wrappers/bind"
	"github.com/alkemics/goflow/wrappers/gonodes"
	"github.com/alkemics/goflow/wrappers/varnames"
)

func TestCycles(t *testing.T) {
	var loader gfgo.NodeLoader
	err := loader.Load("github.com/alkemics/goflow/example/nodes")
	require.NoError(t, err)

	wrappers := []goflow.GraphWrapper{
		gonodes.Wrapper(&loader),
		bind.Wrapper,
		varnames.Wrapper,
		goflow.FromNodeWrapper(after.Wrapper),
		varnames.CompilableWrapper,
	}

	checkers := []goflow.Checker{
		cycles.Check,
	}
	require.NoError(t, gfgo.TestCheck(t, wrappers, checkers, "ok.yml"))
	require.Error(t, gfgo.TestCheck(t, wrappers, checkers, "ko.yml"))
}
