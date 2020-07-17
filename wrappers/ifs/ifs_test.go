package ifs_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alkemics/goflow"
	"github.com/alkemics/goflow/gfutil/gfgo"
	"github.com/alkemics/goflow/wrappers/bind"
	"github.com/alkemics/goflow/wrappers/constants"
	"github.com/alkemics/goflow/wrappers/gonodes"
	"github.com/alkemics/goflow/wrappers/ifs"
	"github.com/alkemics/goflow/wrappers/imports"
	"github.com/alkemics/goflow/wrappers/types"
	"github.com/alkemics/goflow/wrappers/varnames"
)

func TestIfs(t *testing.T) {
	require := require.New(t)

	wd, err := os.Getwd()
	require.NoError(err)
	require.NoError(os.Chdir("../.."))

	var loader gfgo.NodeLoader
	err = loader.Load("github.com/alkemics/goflow/example/nodes")
	require.NoError(err)

	wraps := []goflow.GraphWrapper{
		gonodes.Wrapper(&loader),
		bind.Wrapper,
		goflow.FromNodeWrapper(ifs.Wrapper),
		constants.Wrapper(
			"github.com/alkemics/goflow/example/constants/...",
		),
		varnames.Wrapper,
		types.Wrapper,
		imports.Merger,
		varnames.CompilableWrapper,
	}

	require.NoError(os.Chdir(wd))
	gfgo.TestGenerate(t, wraps, "ifs.yml", gfgo.TestCases{})
}
