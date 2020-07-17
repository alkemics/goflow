package ifs

import "github.com/alkemics/goflow"

func Wrapper(unmarshal func(interface{}) error, node goflow.NodeRenderer) (goflow.NodeRenderer, error) {
	var withIf struct {
		Conditions []string `yaml:"if"`
	}
	if err := unmarshal(&withIf); err != nil {
		return nil, err
	}

	if len(withIf.Conditions) == 0 {
		return node, nil
	}

	return nodeRenderer{
		NodeRenderer: node,
		conditions:   getConditions(withIf.Conditions),
	}, nil
}
