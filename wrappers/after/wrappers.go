package after

import "github.com/alkemics/goflow"

type after struct {
	After []string `yaml:"after"`
}

func Wrapper(unmarshal func(interface{}) error, node goflow.NodeRenderer) (goflow.NodeRenderer, error) {
	var n after
	if err := unmarshal(&n); err != nil {
		return nil, err
	}

	if len(n.After) == 0 {
		return node, nil
	}

	return nodeRenderer{
		NodeRenderer:    node,
		previousNodeIDs: n.After,
	}, nil
}
