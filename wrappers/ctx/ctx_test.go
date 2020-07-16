package ctx_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alkemics/goflow"
	"github.com/alkemics/goflow/gfutil/gfgo"
	"github.com/alkemics/goflow/wrappers/bind"
	"github.com/alkemics/goflow/wrappers/ctx"
	"github.com/alkemics/goflow/wrappers/gonodes"
	"github.com/alkemics/goflow/wrappers/imports"
	"github.com/alkemics/goflow/wrappers/inputs"
	"github.com/alkemics/goflow/wrappers/types"
	"github.com/alkemics/goflow/wrappers/varnames"
)

func TestCtx(t *testing.T) {
	require := require.New(t)

	wd, err := os.Getwd()
	require.NoError(err)
	require.NoError(os.Chdir("../.."))

	var loader gfgo.NodeLoader
	err = loader.Load("github.com/alkemics/goflow/example/nodes")
	require.NoError(err)

	wraps := []goflow.GraphWrapper{
		inputs.Wrapper,
		gonodes.Wrapper(&loader),
		ctx.Wrapper,
		bind.Wrapper,
		varnames.Wrapper,
		types.Wrapper,
		inputs.TypeWrapper,
		imports.Merger,
		varnames.CompilableWrapper,
	}

	require.NoError(os.Chdir(wd))

	testCases := gfgo.TestCases{
		Imports: []string{"context"},
		Tests: []gfgo.TestCase{
			{
				Test: "g.Run(context.Background())",
			},
		},
	}
	gfgo.TestGenerate(t, wraps, "ctx.yml", testCases)
}
