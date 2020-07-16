package gfutil

import (
	"testing"

	"github.com/alkemics/goflow"
	"github.com/stretchr/testify/suite"
)

type DependenciesSuite struct {
	suite.Suite

	dependencyMap map[string][]string
}

func (s *DependenciesSuite) SetupTest() {
	s.dependencyMap = map[string][]string{
		"a": nil,
		"b": {"a"},
		"c": {"a"},
		"d": {"b", "c"},
		"e": nil,
		"f": {"e", "d"},
	}
}

func (s *DependenciesSuite) TestUnravelDependencies() {
	s.dependencyMap["g"] = []string{"u"}
	expected := map[string][]string{
		"a": nil,
		"b": {"a"},
		"c": {"a"},
		"d": {"a", "b", "c"},
		"e": nil,
		"f": {"a", "b", "c", "d", "e"},
		"g": {"u"},
	}
	unraveled := UnravelDependencies(s.dependencyMap)
	s.Assert().Equal(expected, unraveled)
}

type dependencyNode struct {
	id   string
	deps []string
}

func (n dependencyNode) ID() string                                         { return n.id }
func (n dependencyNode) Previous() []string                                 { return n.deps }
func (n dependencyNode) Imports() []goflow.Import                           { return nil }
func (n dependencyNode) Doc() string                                        { return "" }
func (n dependencyNode) Dependencies() []goflow.Field                       { return nil }
func (n dependencyNode) Inputs() []goflow.Field                             { return nil }
func (n dependencyNode) Outputs() []goflow.Field                            { return nil }
func (n dependencyNode) Run(inputs, outputs []goflow.Field) (string, error) { return "", nil }

func (s *DependenciesSuite) TestSortNode() {
	s.dependencyMap["d'"] = []string{"b", "c"}

	nodes := make([]goflow.NodeRenderer, 0)
	for id, deps := range s.dependencyMap {
		nodes = append(nodes, dependencyNode{
			id:   id,
			deps: deps,
		})
	}

	SortNodes(nodes)

	ids := make([]string, len(nodes))
	for i, node := range nodes {
		ids[i] = node.ID()
	}

	expected := []string{
		"a", "e",
		"b", "c",
		"d", "d'",
		"f",
	}
	s.Assert().Equal(expected, ids)
}

func TestDependenciesSuite(t *testing.T) {
	suite.Run(t, new(DependenciesSuite))
}
