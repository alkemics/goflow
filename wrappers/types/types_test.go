package types_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alkemics/goflow"
	"github.com/alkemics/goflow/gfutil/gfgo"
	"github.com/alkemics/goflow/wrappers/after"
	"github.com/alkemics/goflow/wrappers/bind"
	"github.com/alkemics/goflow/wrappers/constants"
	"github.com/alkemics/goflow/wrappers/ctx"
	"github.com/alkemics/goflow/wrappers/gonodes"
	"github.com/alkemics/goflow/wrappers/ifs"
	"github.com/alkemics/goflow/wrappers/imports"
	"github.com/alkemics/goflow/wrappers/inputs"
	"github.com/alkemics/goflow/wrappers/outputs"
	"github.com/alkemics/goflow/wrappers/types"
	"github.com/alkemics/goflow/wrappers/varnames"
)

func TestTypes(t *testing.T) {
	require := require.New(t)

	wd, err := os.Getwd()
	require.NoError(err)
	require.NoError(os.Chdir("../.."))

	var loader gfgo.NodeLoader
	err = loader.Load("github.com/alkemics/goflow/example/nodes")
	require.NoError(err)

	wraps := []goflow.GraphWrapper{
		inputs.Wrapper,
		goflow.FromNodeWrapper(after.Wrapper),
		gonodes.Wrapper(&loader),
		gonodes.DepWrapper,
		ctx.Wrapper,
		bind.Wrapper,
		goflow.FromNodeWrapper(ifs.Wrapper),
		constants.Wrapper(
			"github.com/alkemics/goflow/example/constants...",
		),
		outputs.Wrapper,
		varnames.Wrapper,
		types.Wrapper, // 12
		inputs.TypeWrapper,
		outputs.NameWrapper,
		varnames.CompilableWrapper,
		imports.Merger,
	}

	require.NoError(os.Chdir(wd))

	testCases := gfgo.TestCases{
		Imports: []string{
			"context",
		},
		Tests: []gfgo.TestCase{
			{Test: `
			add, add10, forwardedInputs := g.Run(context.Background(), 1, []uint{2, 3}, 4, false)
			if add != 28 { panic("add should be 28")}
			if add10 != 11 { panic("add10 should be 11")}
			if fmt.Sprint(forwardedInputs) != "[1 2 3 4]" { panic(fmt.Sprintf("forwardedInputs should be [1 2 3 4], got %v", forwardedInputs)) }
`},
		},
	}
	gfgo.TestGenerate(t, wraps, "types.yml", testCases)
}
