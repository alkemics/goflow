package gfutil

import "github.com/alkemics/goflow"

type DummyNodeRenderer struct {
	IDVal           string
	PreviousVal     []string
	ImportsVal      []goflow.Import
	DocVal          string
	DependenciesVal []goflow.Field
	InputsVal       []goflow.Field
	OutputsVal      []goflow.Field
	RunFunc         func(inputs, outputs []goflow.Field) (string, error)
}

func (d DummyNodeRenderer) ID() string {
	return d.IDVal
}

func (d DummyNodeRenderer) Previous() []string {
	return d.PreviousVal
}

func (d DummyNodeRenderer) Imports() []goflow.Import {
	return d.ImportsVal
}

func (d DummyNodeRenderer) Doc() string {
	return d.DocVal
}

func (d DummyNodeRenderer) Dependencies() []goflow.Field {
	return d.DependenciesVal
}

func (d DummyNodeRenderer) Inputs() []goflow.Field {
	return d.InputsVal
}

func (d DummyNodeRenderer) Outputs() []goflow.Field {
	return d.OutputsVal
}

func (d DummyNodeRenderer) Run(inputs, outputs []goflow.Field) (string, error) {
	return d.RunFunc(inputs, outputs)
}

var _ goflow.NodeRenderer = DummyNodeRenderer{}

func StringRunFunc(str string) func(inputs, outputs []goflow.Field) (string, error) {
	return func(inputs, outputs []goflow.Field) (string, error) {
		return str, nil
	}
}

func ErrRunFunc(err error) func(inputs, outputs []goflow.Field) (string, error) {
	return func(inputs, outputs []goflow.Field) (string, error) {
		return "", err
	}
}
