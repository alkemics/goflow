package gfgo

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type NodeLoaderSuite struct {
	suite.Suite
}

func (s NodeLoaderSuite) TestLoadNodes_ignoredNodes() {
	nl := NodeLoader{}
	s.Assert().NoError(nl.Load("github.com/alkemics/goflow/gfutil/gfgo/internal/nodes"))

	var node Node

	node = nl.Find("nodes.IgnoreType")
	s.Assert().Zero(node)

	node = nl.Find("nodes.IgnoredFunction")
	s.Assert().Zero(node)

	node = nl.Find("nodes.IgnoredMethod.Delete")
	s.Assert().Zero(node)

	node = nl.Find("nodes.IgnoredMethod.Create")
	s.Assert().NotZero(node)
}

func (s NodeLoaderSuite) TestNodeLoader() {
	nl := NodeLoader{}
	s.Assert().NoError(nl.Load("github.com/alkemics/goflow/gfutil/gfgo/internal/nodes"))

	var node Node

	node = nl.Find("nodes.Adder")
	s.Assert().NotZero(node)

	nl.nodes = nil
	node = nl.Find("nodes.Adder")
	s.Assert().Zero(node)

	s.Assert().NoError(nl.Refresh("nodes"))
	node = nl.Find("nodes.Adder")
	s.Assert().NotZero(node)
}

func TestNodeLoader(t *testing.T) {
	suite.Run(t, new(NodeLoaderSuite))
}
