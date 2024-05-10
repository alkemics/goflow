package gfgo

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
	"text/template"

	"github.com/alkemics/goflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const mainTmpl = `
// Code generated by lib-go/goflow DO NOT EDIT.

// +build !codeanalysis

package main

import (
	"fmt"
	{{ range .Imports }}
	"{{ . }}"
	{{ end }}
)

func assert(b bool) {
	if !b {
		panic("assertion failed")
	}
}

func main() {
	var g {{ .Name }}
	{{ range .Tests }}
	g = New{{ $.Name }}({{ .Deps }})
	{{ .Test }}
	fmt.Println("Test done.")
	{{ end }}
	fmt.Println("All tests done.")
}
`

type TestCases struct {
	Imports []string
	Tests   []TestCase
}

type TestCase struct {
	Deps string
	Test string
}

// TestGenerate should be used to simplify testing the wrappers.
//
// Each element of testCases should be valid go code. The graph
// created is called g.
func TestGenerate(t *testing.T, wrappers []goflow.GraphWrapper, filename string, testCases TestCases) {
	require := require.New(t)

	graph, err := goflow.Load(filename, wrappers)
	require.NoError(err)

	goFilename := path.Join("graph", strings.ReplaceAll(filename, ".yml", ".go"))
	buf := bytes.Buffer{}
	w := NewWriter(
		&buf,
		goFilename,
	)
	err = Generate(w, graph)
	require.NoError(err)
	require.NoError(w.Flush())

	if _, err := os.Open("graph"); os.IsNotExist(err) {
		require.NoError(os.Mkdir("graph", 0o755))
	} else if err != nil {
		require.NoError(err)
	}
	err = os.WriteFile(goFilename, buf.Bytes(), 0o600)
	require.NoError(err)

	mainFilename := path.Join("graph", "main.go")
	tmpl, err := template.New("main").Parse(mainTmpl)
	require.NoError(err)

	if len(testCases.Tests) == 0 {
		testCases = TestCases{
			Tests: []TestCase{
				{
					Deps: "",
					Test: "g.Run()",
				},
			},
		}
	}

	buf = bytes.Buffer{}
	w = NewWriter(
		&buf,
		mainFilename,
	)
	err = tmpl.Execute(w, struct {
		Name    string
		Imports []string
		Tests   []TestCase
	}{
		Name:    graph.Name(),
		Imports: testCases.Imports,
		Tests:   testCases.Tests,
	})
	require.NoError(err)
	require.NoError(w.Flush())

	require.NoError(err)
	require.NoError(os.WriteFile(mainFilename, buf.Bytes(), 0o600))

	cmd := exec.Command("go", "run", "./graph")
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if !assert.NoError(t, cmd.Run()) {
		t.Log(errOut.String())
	}
	t.Logf("Graph output:\n%s", out.String())
}

func TestLint(t *testing.T, linters []goflow.Linter, filename string) error {
	wrappers := make([]goflow.GraphWrapper, len(linters))
	for i, linter := range linters {
		wrappers[i] = goflow.FromLinter(linter)
	}

	_, err := goflow.Load(filename, wrappers)
	return err
}

func TestCheck(t *testing.T, wrappers []goflow.GraphWrapper, checkers []goflow.Checker, filename string) error {
	for _, checker := range checkers {
		wrappers = append(wrappers, goflow.FromChecker(checker))
	}

	_, err := goflow.Load(filename, wrappers)
	return err
}
