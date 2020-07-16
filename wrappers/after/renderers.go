package after

import "github.com/alkemics/goflow"

type nodeRenderer struct {
	goflow.NodeRenderer

	previousNodeIDs []string
}

func (n nodeRenderer) Previous() []string {
	return append(n.NodeRenderer.Previous(), n.previousNodeIDs...)
}
