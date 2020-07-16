package imports

import (
	"path"

	"github.com/alkemics/goflow"
)

func Merger(_ func(interface{}) error, graph goflow.GraphRenderer) (goflow.GraphRenderer, error) {
	return merger{GraphRenderer: graph}, nil
}

type merger struct {
	goflow.GraphRenderer
}

func (g merger) Imports() []goflow.Import {
	imports := Merge(g.GraphRenderer.Imports())

	for _, node := range g.Nodes() {
		imports = Merge(imports, node.Imports()...)
	}

	// Remove the package if redundant.
	for i, imp := range imports {
		if path.Base(imp.Dir) == imp.Pkg {
			imports[i].Pkg = ""
		}
	}

	return imports
}

// Compile time check
var _ goflow.GraphRenderer = merger{}

func Merge(imports []goflow.Import, others ...goflow.Import) []goflow.Import {
	merged := make([]goflow.Import, 0)
	importDirs := make(map[string]struct{})

	for _, imp := range imports {
		if _, ok := importDirs[imp.Dir]; !ok {
			importDirs[imp.Dir] = struct{}{}

			merged = append(merged, imp)
		}
	}

	for _, imp := range others {
		if _, ok := importDirs[imp.Dir]; !ok {
			importDirs[imp.Dir] = struct{}{}

			merged = append(merged, imp)
		}
	}

	return merged
}
