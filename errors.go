package goflow

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Error struct {
	Filename string

	Err error
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Filename, e.Err.Error())
}

func (e Error) Unwrap() error {
	return e.Err
}

func (e Error) Format(s fmt.State, verb rune) {
	if verb != 'v' || (!s.Flag('+')) {
		fmt.Fprint(s, e.Error())
		return
	}

	errMsg := fmt.Sprintf("  %v", e.Err)

	var multiErr MultiError
	if errors.As(e.Err, &multiErr) {
		flat := multiErr.flatten()
		strs := make([]string, len(flat.Errs))
		for i, err := range flat.Errs {
			strs[i] = fmt.Sprintf("  %v", err)
		}
		errMsg = strings.Join(strs, "\n")
	}

	fmt.Fprintf(s, `%s:
%v
`,
		e.Filename,
		errMsg,
	)
}

type MultiError struct {
	Errs []error
}

func (e MultiError) Error() string {
	if len(e.Errs) == 1 {
		return e.Errs[0].Error()
	}

	strs := make([]string, len(e.Errs))
	for i, err := range e.Errs {
		strs[i] = err.Error()
	}
	return fmt.Sprintf("[%s]", strings.Join(strs, ", "))
}

func (e MultiError) Is(err error) bool {
	for _, sub := range e.Errs {
		if errors.Is(sub, err) {
			return true
		}
	}
	return false
}

func (e MultiError) As(target interface{}) bool {
	for _, sub := range e.Errs {
		if errors.As(sub, target) {
			return true
		}
	}
	return false
}

func (e MultiError) flatten() MultiError {
	errs := make([]error, 0, len(e.Errs))
	for _, err := range e.Errs {
		var me MultiError
		if errors.As(err, &me) {
			errs = append(errs, me.flatten().Errs...)
		} else {
			errs = append(errs, err)
		}
	}
	return MultiError{Errs: errs}
}

type GraphError struct {
	Wrapper string
	Err     error
}

func (e GraphError) Error() string {
	return fmt.Sprintf("%v (%s)", e.Err, e.Wrapper)
}

func (e GraphError) Unwrap() error {
	return e.Err
}

type NodeError struct {
	ID      string
	Wrapper string

	Err error
}

func (e NodeError) Error() string {
	return fmt.Sprintf("%s: %v (%s)", e.ID, e.Err, e.Wrapper)
}

func (e NodeError) Unwrap() error {
	return e.Err
}

var yamlErrorRegex = regexp.MustCompile(`(?s)(?:yaml: )?(?:unmarshal errors:.)?line (\d+): (.*)`)

type YAMLError struct {
	Line int
	Err  error
}

func (e YAMLError) Error() string {
	return fmt.Sprintf("line %d: %v (yaml)", e.Line, e.Err.Error())
}

func (e YAMLError) Unwrap() error {
	return e.Err
}

func ParseYAMLError(err error) error {
	if err == nil {
		return nil
	}

	match := yamlErrorRegex.FindStringSubmatch(err.Error())
	if len(match) == 0 {
		return err
	}

	line, _ := strconv.Atoi(match[1])
	return YAMLError{
		Line: line,
		Err:  errors.New(strings.TrimSpace(match[2])),
	}
}
