package goflow

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
)

type ErrorSuite struct {
	suite.Suite
}

func TestErrorSuite(t *testing.T) {
	suite.Run(t, new(ErrorSuite))
}

func (s *ErrorSuite) TestMultiError_flatten() {
	merr := MultiError{
		Errs: []error{
			errors.New("1"),
			errors.New("2"),
			MultiError{
				Errs: []error{
					errors.New("3"),
					errors.New("4"),
				},
			},
			errors.New("5"),
			MultiError{
				Errs: []error{
					errors.New("6"),
					MultiError{
						Errs: []error{
							errors.New("7"),
							errors.New("8"),
						},
					},
					MultiError{
						Errs: []error{
							errors.New("9"),
							errors.New("10"),
						},
					},
				},
			},
		},
	}

	flattened := merr.flatten()
	expected := MultiError{
		Errs: []error{
			errors.New("1"),
			errors.New("2"),
			errors.New("3"),
			errors.New("4"),
			errors.New("5"),
			errors.New("6"),
			errors.New("7"),
			errors.New("8"),
			errors.New("9"),
			errors.New("10"),
		},
	}
	s.Assert().Equal(expected, flattened)
}

func (s *ErrorSuite) TestParseYAMLError_nil() {
	yamlErr := ParseYAMLError(nil)
	s.Assert().Nil(yamlErr)
}

func (s *ErrorSuite) TestParseYAMLError_notYAML() {
	err := errors.New("not a yaml error")
	yamlErr := ParseYAMLError(err)
	s.Assert().Equal(err, yamlErr)
}

func (s *ErrorSuite) TestParseYAMLError() {
	r := strings.NewReader(`
- key: value
  error
`)
	err := yaml.NewDecoder(r).Decode(nil)
	yamlErr := ParseYAMLError(err)
	expected := YAMLError{
		Line: 4,
		Err:  errors.New("could not find expected ':'"),
	}
	s.Assert().Equal(expected, yamlErr)
}
