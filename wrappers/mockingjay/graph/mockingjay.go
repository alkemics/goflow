// Code generated by goflow DO NOT EDIT.

//go:build !codeanalysis
// +build !codeanalysis

package main

import (
	"context"

	"github.com/alkemics/goflow/example/nodes"
	"github.com/alkemics/goflow/wrappers/mockingjay"
)

/*

 */
type Mockingjay struct{}

func NewMockingjay() Mockingjay {
	return Mockingjay{}
}

func newMockingjay(id string) Mockingjay {
	return Mockingjay{}
}

/*

 */
func (g *Mockingjay) Run(ctx context.Context, a int, b int) (sum int) {

	// __add_a outputs
	var __add_a_aggregated int

	// __add_b outputs
	var __add_b_aggregated int

	// __ctx outputs
	var __ctx_ctx context.Context

	// __output_sum_builder outputs
	var __output_sum_builder_sum int

	// __print_values outputs
	var __print_values_aggregated []interface{}

	// add outputs
	var add_sum int

	// inputs outputs
	var inputs_a int
	var inputs_b int

	// print outputs

	igniteNodeID := "ignite"
	doneNodeID := "done"

	done := make(chan string)
	defer close(done)

	steps := map[string]struct {
		deps        map[string]struct{}
		run         func()
		alreadyDone bool
	}{

		"__add_a": {
			deps: map[string]struct{}{
				"inputs":     {},
				igniteNodeID: {},
			},
			run: func() {
				__add_a_aggregated = inputs_a
				done <- "__add_a"
			},
			alreadyDone: false,
		},
		"__add_b": {
			deps: map[string]struct{}{
				"inputs":     {},
				igniteNodeID: {},
			},
			run: func() {
				__add_b_aggregated = inputs_b
				done <- "__add_b"
			},
			alreadyDone: false,
		},
		"__ctx": {
			deps: map[string]struct{}{
				igniteNodeID: {},
			},
			run: func() {
				__ctx_ctx = ctx
				done <- "__ctx"
			},
			alreadyDone: false,
		},
		"__output_sum_builder": {
			deps: map[string]struct{}{
				"add":        {},
				igniteNodeID: {},
			},
			run: func() {

				__output_sum_builder_sum = add_sum
				sum = __output_sum_builder_sum
				done <- "__output_sum_builder"
			},
			alreadyDone: false,
		},
		"__print_values": {
			deps: map[string]struct{}{
				"add":        {},
				igniteNodeID: {},
			},
			run: func() {
				__print_values_aggregated = append(__print_values_aggregated, "sum")
				__print_values_aggregated = append(__print_values_aggregated, add_sum)
				done <- "__print_values"
			},
			alreadyDone: false,
		},
		"add": {
			deps: map[string]struct{}{
				"__ctx":      {},
				"__add_a":    {},
				"__add_b":    {},
				igniteNodeID: {},
			},
			run: func() {

				var _mock []interface{}
				if _mocks, ok := ctx.Value(mockingjay.ContextKey).(map[string][]interface{}); ok && _mocks != nil {
					m, ok := _mocks["add"]
					if ok {
						_mock = m
					}
				}

				if _mock != nil {
					add_sum = _mock[0].(int)
				} else {
					add_sum = nodes.Adder(__add_a_aggregated, __add_b_aggregated)
				}

				done <- "add"
			},
			alreadyDone: false,
		},
		"inputs": {
			deps: map[string]struct{}{
				igniteNodeID: {},
			},
			run: func() {

				var _mock []interface{}
				if _mocks, ok := ctx.Value(mockingjay.ContextKey).(map[string][]interface{}); ok && _mocks != nil {
					m, ok := _mocks["inputs"]
					if ok {
						_mock = m
					}
				}

				if _mock != nil {
					inputs_a = _mock[0].(int)
					inputs_b = _mock[1].(int)
				} else {
					inputs_a = a
					inputs_b = b
				}

				done <- "inputs"
			},
			alreadyDone: false,
		},
		"print": {
			deps: map[string]struct{}{
				"__ctx":          {},
				"__print_values": {},
				igniteNodeID:     {},
			},
			run: func() {

				var _mock []interface{}
				if _mocks, ok := ctx.Value(mockingjay.ContextKey).(map[string][]interface{}); ok && _mocks != nil {
					m, ok := _mocks["print"]
					if ok {
						_mock = m
					}
				}

				if _mock != nil {

				} else {
					nodes.PrinterCtx(__ctx_ctx, __print_values_aggregated)
				}

				done <- "print"
			},
			alreadyDone: false,
		},
		igniteNodeID: {
			deps: map[string]struct{}{},
			run: func() {
				done <- igniteNodeID
			},
			alreadyDone: false,
		},
		doneNodeID: {
			deps: map[string]struct{}{
				"__add_a":              {},
				"__add_b":              {},
				"__ctx":                {},
				"__output_sum_builder": {},
				"__print_values":       {},
				"add":                  {},
				"inputs":               {},
				"print":                {},
			},
			run: func() {
				done <- doneNodeID
			},
			alreadyDone: false,
		},
	}

	// Ignite
	ignite := steps[igniteNodeID]
	ignite.alreadyDone = true
	steps[igniteNodeID] = ignite
	go steps[igniteNodeID].run()

	// Resolve the graph
	for resolved := range done {
		if resolved == doneNodeID {
			// If all the graph was resolved, get out of the loop
			break
		}

		for name, step := range steps {
			delete(step.deps, resolved)
			if len(step.deps) == 0 && !step.alreadyDone {
				step.alreadyDone = true
				steps[name] = step
				go step.run()
			}
		}
	}

	return sum
}
