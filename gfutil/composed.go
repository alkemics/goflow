package gfutil

import "github.com/alkemics/goflow"

func ComposeWrappers(wrappers ...goflow.GraphWrapper) goflow.GraphWrapper {
	return func(unmarshal func(interface{}) error, graph goflow.GraphRenderer) (goflow.GraphRenderer, error) {
		for _, wrapper := range wrappers {
			g, err := wrapper(unmarshal, graph)
			if err != nil {
				return nil, err
			}

			graph = g
		}

		return graph, nil
	}
}
