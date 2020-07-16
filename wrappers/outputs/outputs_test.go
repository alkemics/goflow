package outputs_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alkemics/goflow"
	"github.com/alkemics/goflow/gfutil/gfgo"
	"github.com/alkemics/goflow/wrappers/gonodes"
	"github.com/alkemics/goflow/wrappers/imports"
	"github.com/alkemics/goflow/wrappers/outputs"
	"github.com/alkemics/goflow/wrappers/types"
	"github.com/alkemics/goflow/wrappers/varnames"
)

func TestOutputs(t *testing.T) {
	require := require.New(t)

	wd, err := os.Getwd()
	require.NoError(err)
	require.NoError(os.Chdir("../.."))

	var loader gfgo.NodeLoader
	err = loader.Load("github.com/alkemics/goflow/example/nodes")
	require.NoError(err)

	wraps := []goflow.GraphWrapper{
		gonodes.Wrapper(&loader),
		outputs.Wrapper,
		varnames.Wrapper,
		types.Wrapper,
		imports.Merger,
		outputs.NameWrapper,
		varnames.CompilableWrapper,
	}

	require.NoError(os.Chdir(wd))

	testCases := gfgo.TestCases{
		Tests: []gfgo.TestCase{
			{Test: `
ri, ris := g.Run()
if ri < 0 || ri > 100 {
	panic("invalid ri")
}
if len(ris) != 2 {
	panic("ris should have length 2")
}
			`},
		},
	}

	gfgo.TestGenerate(t, wraps, "outputs.yml", testCases)
}
