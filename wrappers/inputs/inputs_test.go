package inputs_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alkemics/goflow"
	"github.com/alkemics/goflow/gfutil/gfgo"
	"github.com/alkemics/goflow/wrappers/bind"
	"github.com/alkemics/goflow/wrappers/constants"
	"github.com/alkemics/goflow/wrappers/ctx"
	"github.com/alkemics/goflow/wrappers/gonodes"
	"github.com/alkemics/goflow/wrappers/ifs"
	"github.com/alkemics/goflow/wrappers/imports"
	"github.com/alkemics/goflow/wrappers/inputs"
	"github.com/alkemics/goflow/wrappers/types"
	"github.com/alkemics/goflow/wrappers/varnames"
)

func TestInputs(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir("../.."))

	var loader gfgo.NodeLoader
	err = loader.Load("github.com/alkemics/goflow/example/nodes")
	require.NoError(t, err)

	wraps := []goflow.GraphWrapper{
		inputs.Wrapper,
		gonodes.Wrapper(&loader),
		ctx.Wrapper,
		bind.Wrapper,
		goflow.FromNodeWrapper(ifs.Wrapper),
		constants.Wrapper(
			"github.com/alkemics/goflow/example/constants/...",
		),
		varnames.Wrapper,
		types.Wrapper,
		inputs.TypeWrapper,
		imports.Merger,
		varnames.CompilableWrapper,
	}

	require.NoError(t, os.Chdir(wd))

	testCases := gfgo.TestCases{
		Imports: []string{
			"context",
		},
		Tests: []gfgo.TestCase{
			{Test: "g.Run(context.Background(), 1, []int{2, 3}, false)"},
		},
	}
	gfgo.TestGenerate(t, wraps, "inputs.yml", testCases)
}
