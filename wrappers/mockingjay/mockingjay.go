package mockingjay

import (
	"context"
	"fmt"
	"strings"

	"github.com/alkemics/goflow"
	"github.com/alkemics/goflow/wrappers/ctx"
)

type ContextKeyType string

const ContextKey ContextKeyType = "mockingjay"

func WithMock(ctx context.Context, nodeID string, values ...interface{}) context.Context {
	m, ok := ctx.Value(ContextKey).(map[string][]interface{})
	if !ok || m == nil {
		m = make(map[string][]interface{})
	}

	m[nodeID] = values
	return context.WithValue(ctx, ContextKey, m)
}

// Mock is a goflow.NodeWrapper that lets you mock nodes at run time.
//
// The mock will extract the value from the context at ContextKey as a
// map[string][]interface{}, looking in the map via the node id to find
// the mocked returns. If the node is not mocked, it is executed normally.
//
// The WithMock helper function is provided to make it easier to mock  node:
//
//	ctx := mockingjay.WithMock(ctx, "myNode", 42)
//	graph.Run(ctx, ...)
func Mock(_ func(interface{}) error, node goflow.NodeRenderer) (goflow.NodeRenderer, error) {
	return mocker{
		NodeRenderer: ctx.Injector{
			NodeRenderer: node,
		},
	}, nil
}

type mocker struct {
	goflow.NodeRenderer
}

func (m mocker) Run(inputs, outputs []goflow.Field) (string, error) {
	sub, err := m.NodeRenderer.Run(inputs, outputs)
	if err != nil {
		return "", err
	}

	mocked := make([]string, len(outputs))
	for i, output := range outputs {
		mocked[i] = fmt.Sprintf("%s = _mock[%d].(%s)", output.Name, i, output.Type)
	}

	return fmt.Sprintf(`
var _mock []interface{}
if _mocks, ok := ctx.Value(mockingjay.ContextKey).(map[string][]interface{}); ok && _mocks != nil {
	m, ok := _mocks["%s"]
	if ok {
		_mock = m
	}
}

if _mock != nil {
	%s
} else {
	%s
}
`,
		m.ID(),
		strings.Join(mocked, "\n"),
		sub,
	), nil
}
