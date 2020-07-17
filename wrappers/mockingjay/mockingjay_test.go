package mockingjay_test

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
	"github.com/alkemics/goflow/wrappers/imports"
	"github.com/alkemics/goflow/wrappers/inputs"
	"github.com/alkemics/goflow/wrappers/mockingjay"
	"github.com/alkemics/goflow/wrappers/outputs"
	"github.com/alkemics/goflow/wrappers/types"
	"github.com/alkemics/goflow/wrappers/varnames"
)

func TestVia(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir("../.."))

	var loader gfgo.NodeLoader
	err = loader.Load("github.com/alkemics/goflow/example/nodes")
	require.NoError(t, err)

	wraps := []goflow.GraphWrapper{
		inputs.Wrapper,
		gonodes.Wrapper(&loader),
		gonodes.DepWrapper,
		goflow.FromNodeWrapper(mockingjay.Mock),
		ctx.Wrapper,
		bind.Wrapper,
		constants.Wrapper(
			"github.com/alkemics/goflow/example/constants/...",
		),
		outputs.Wrapper,
		varnames.Wrapper,
		types.Wrapper,
		inputs.TypeWrapper,
		imports.Merger,
		outputs.NameWrapper,
		varnames.CompilableWrapper,
	}

	require.NoError(t, os.Chdir(wd))

	testCases := gfgo.TestCases{
		Imports: []string{
			"context",
			"github.com/alkemics/goflow/wrappers/mockingjay",
		},
		Tests: []gfgo.TestCase{
			{
				Test: `
ctx := mockingjay.WithMock(context.Background(), "add", 1)
sum := g.Run(ctx, 10, 20)
if sum != 1 {
	panic(fmt.Sprintf("expected 1 got %v\n", sum))
}
`,
			},
		},
	}
	gfgo.TestGenerate(t, wraps, "mockingjay.yml", testCases)
}
