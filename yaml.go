package goflow

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Load will decode the yaml file at yamlFilename into a GraphRenderer, installing all
// the wrappers.
//
// Even if a wrapper returns an error, this function will go through all the wrappers.
func Load(yamlFilename string, wrappers []GraphWrapper) (GraphRenderer, error) {
	file, err := os.Open(yamlFilename)
	if err != nil {
		return nil, Error{
			Filename: yamlFilename,
			Err:      err,
		}
	}
	defer file.Close()

	var g graphLoader
	if err := yaml.NewDecoder(file).Decode(&g); err != nil {
		return nil, ParseYAMLError(err)
	}

	if err := file.Close(); err != nil {
		return nil, Error{
			Filename: yamlFilename,
			Err:      err,
		}
	}

	graph, err := loadGraph(g, wrappers)
	if err != nil {
		return nil, Error{
			Filename: yamlFilename,
			Err:      err,
		}
	}

	return graph, nil
}

func loadGraph(g graphLoader, wrappers []GraphWrapper) (GraphRenderer, error) {
	var wrapped GraphRenderer = g

	errs := make([]error, 0)
	for _, w := range wrappers {
		w, err := w(g.Unmarshal, wrapped)
		if err != nil {
			errs = append(errs, err)
		}
		if w != nil {
			wrapped = w
		}
	}

	if len(errs) > 0 {
		return nil, MultiError{Errs: errs}
	}

	return wrapped, nil
}

type graphLoader struct {
	NameVal   string
	PkgVal    string
	NodesVal  []nodeLoader
	Unmarshal func(interface{}) error
}

func (g graphLoader) Name() string          { return g.NameVal }
func (g graphLoader) Pkg() string           { return g.PkgVal }
func (g graphLoader) Imports() []Import     { return nil }
func (g graphLoader) Doc() string           { return "" }
func (g graphLoader) Dependencies() []Field { return nil }
func (g graphLoader) Inputs() []Field       { return nil }
func (g graphLoader) Outputs() []Field      { return nil }
func (g graphLoader) Nodes() []NodeRenderer {
	nodes := make([]NodeRenderer, len(g.NodesVal))
	for i, n := range g.NodesVal {
		nodes[i] = n
	}
	return nodes
}

func (g *graphLoader) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var gl struct {
		Name  string       `yaml:"name"`
		Pkg   string       `yaml:"package"`
		Nodes []nodeLoader `yaml:"nodes"`
	}
	if err := unmarshal(&gl); err != nil {
		return err
	}
	g.NameVal = gl.Name
	g.PkgVal = gl.Pkg
	if g.PkgVal == "" {
		g.PkgVal = "graphs"
	}
	g.NodesVal = gl.Nodes
	g.Unmarshal = func(v interface{}) error {
		return ParseYAMLError(unmarshal(v))
	}
	return nil
}

type nodeLoader struct {
	IDVal     string
	Unmarshal func(interface{}) error
}

func (n nodeLoader) ID() string                                  { return n.IDVal }
func (n nodeLoader) Doc() string                                 { return "" }
func (n nodeLoader) Previous() []string                          { return nil }
func (n nodeLoader) Imports() []Import                           { return nil }
func (n nodeLoader) Dependencies() []Field                       { return nil }
func (n nodeLoader) Inputs() []Field                             { return nil }
func (n nodeLoader) Outputs() []Field                            { return nil }
func (n nodeLoader) Run(inputs, outputs []Field) (string, error) { return "", nil }

func (n nodeLoader) UnmarshalFunc() func(interface{}) error {
	return n.Unmarshal
}

func (n *nodeLoader) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var nl struct {
		ID string `yaml:"id"`
	}
	if err := unmarshal(&nl); err != nil {
		return err
	}
	n.IDVal = nl.ID
	n.Unmarshal = func(v interface{}) error {
		return ParseYAMLError(unmarshal(v))
	}
	return nil
}
