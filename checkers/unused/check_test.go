package unused_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alkemics/goflow"
	"github.com/alkemics/goflow/checkers/unused"
	"github.com/alkemics/goflow/gfutil/gfgo"
	"github.com/alkemics/goflow/wrappers/after"
	"github.com/alkemics/goflow/wrappers/bind"
	"github.com/alkemics/goflow/wrappers/gonodes"
	"github.com/alkemics/goflow/wrappers/outputs"
	"github.com/alkemics/goflow/wrappers/varnames"
)

func TestUnused(t *testing.T) {
	require := require.New(t)

	var loader gfgo.NodeLoader
	err := loader.Load("github.com/alkemics/goflow/example/nodes")
	require.NoError(err)

	wrappers := []goflow.GraphWrapper{
		gonodes.Wrapper(&loader),
		bind.Wrapper,
		outputs.Wrapper,
		varnames.Wrapper,
		outputs.NameWrapper,
		goflow.FromNodeWrapper(after.Wrapper),
		varnames.CompilableWrapper,
	}
	checkers := []goflow.Checker{
		unused.Check,
	}
	require.NoError(gfgo.TestCheck(t, wrappers, checkers, "ok.yml"))
	require.Error(gfgo.TestCheck(t, wrappers, checkers, "ko.yml"))
}
