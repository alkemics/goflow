package ids_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alkemics/goflow"
	"github.com/alkemics/goflow/gfutil/gfgo"
	"github.com/alkemics/goflow/linters/ids"
)

func TestIDs(t *testing.T) {
	linters := []goflow.Linter{
		ids.Lint,
	}
	require.NoError(t, gfgo.TestLint(t, linters, "ok.yml"))
	require.Error(t, gfgo.TestLint(t, linters, "ko.yml"))
}
