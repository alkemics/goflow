package after

import (
	"testing"

	"github.com/alkemics/goflow/gfutil"
	"github.com/stretchr/testify/suite"
)

type afterSuite struct {
	suite.Suite
}

func (s afterSuite) TestWrap() {
	node := gfutil.DummyNodeRenderer{
		PreviousVal: []string{
			"a",
		},
	}
	unmarshal := func(v interface{}) error {
		a := v.(*after)
		*a = after{
			After: []string{"b", "c"},
		}

		return nil
	}

	wrapped, err := Wrapper(unmarshal, node)
	s.Assert().NoError(err)
	s.Assert().Equal([]string{"a", "b", "c"}, wrapped.Previous())
}

func (s afterSuite) TestNoAfter() {
	node := gfutil.DummyNodeRenderer{
		PreviousVal: []string{
			"a",
		},
	}
	unmarshal := func(v interface{}) error {
		return nil
	}

	wrapped, err := Wrapper(unmarshal, node)
	s.Assert().NoError(err)
	s.Assert().Equal([]string{"a"}, wrapped.Previous())
}

func TestAfterSuite(t *testing.T) {
	suite.Run(t, new(afterSuite))
}
