package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/alkemics/goflow"
	"github.com/alkemics/goflow/checkers/cycles"
	"github.com/alkemics/goflow/checkers/previous"
	"github.com/alkemics/goflow/checkers/unused"
	"github.com/alkemics/goflow/gfutil"
	"github.com/alkemics/goflow/gfutil/gfgo"
	"github.com/alkemics/goflow/linters/ids"
	"github.com/alkemics/goflow/wrappers/after"
	"github.com/alkemics/goflow/wrappers/bind"
	"github.com/alkemics/goflow/wrappers/constants"
	"github.com/alkemics/goflow/wrappers/ctx"
	"github.com/alkemics/goflow/wrappers/gonodes"
	"github.com/alkemics/goflow/wrappers/ifs"
	"github.com/alkemics/goflow/wrappers/imports"
	"github.com/alkemics/goflow/wrappers/inputs"
	"github.com/alkemics/goflow/wrappers/outputs"
	"github.com/alkemics/goflow/wrappers/types"
	"github.com/alkemics/goflow/wrappers/varnames"
)

func main() {
	var graphNames []string
	if len(os.Args) > 1 {
		for _, gn := range os.Args[1:] {
			if gn != "" {
				graphNames = append(graphNames, gn)
			}
		}
	}

	if len(graphNames) == 0 {
		gns, err := gfutil.FindGraphFileNames("graphs")
		handleErr(err)

		for _, gn := range gns {
			if !strings.Contains(gn, "/errored/") {
				graphNames = append(graphNames, gn)
			}
		}
	}

	var nodeLoader gfgo.NodeLoader

	// This one will fail fast if there is a compilation error in the nodes
	err := nodeLoader.Load(
		path.Join("github.com/alkemics/goflow/example", "nodes", "..."),
		path.Join("github.com/alkemics/goflow/example", "graphs", "..."),
	)
	handleErr(err)

	goNodes := nodeLoader.All()

	genErrs := make([]error, 0, len(graphNames))
	wrappers := getWrappers(&nodeLoader)
	for _, graphName := range graphNames {
		graph, err := goflow.Load(graphName, wrappers)
		if notFoundErr := (gonodes.NotFoundError{}); errors.As(err, &notFoundErr) {
			split := strings.SplitN(notFoundErr.Type, ".", 2)
			if len(split) != 2 {
				genErrs = append(genErrs, err)
				continue
			}

			pkgName := split[0]
			if err := nodeLoader.Refresh(pkgName); err != nil {
				genErrs = append(genErrs, err)
				continue
			}

			goNodes = nodeLoader.All()
			wrappers = getWrappers(&nodeLoader)
			graph, err = goflow.Load(graphName, wrappers)
		}

		if err != nil {
			genErrs = append(genErrs, err)
			continue
		}

		generatedFilename := strings.ReplaceAll(graphName, ".yml", ".go")
		generate := gfgo.WithJSONMarshal(
			gfgo.Generate,
			graphName,
			goNodes,
		)
		buf := bytes.Buffer{}
		w := gfgo.NewWriter(
			gfgo.DebugWriter{Writer: &buf},
			generatedFilename,
		)
		if err := generate(w, graph); err != nil {
			genErrs = append(genErrs, err)
			continue
		}

		if err := w.Flush(); err != nil {
			genErrs = append(genErrs, err)
			continue
		}

		if err := ioutil.WriteFile(generatedFilename, buf.Bytes(), 0o600); err != nil {
			genErrs = append(genErrs, err)
		}

		fmt.Println(generatedFilename, "handled")
	}

	handleErr(genErrs...)

	err = nodeLoader.Load(path.Join("github.com/alkemics/goflow/example", "graphs", "..."))
	handleErr(err)

	buf := bytes.Buffer{}
	nodes := nodeLoader.All()
	sort.SliceStable(nodes, func(i, j int) bool {
		n1 := fmt.Sprintf("%s.%s.%s", nodes[i].Pkg, nodes[i].Typ, nodes[i].Method)
		n2 := fmt.Sprintf("%s.%s.%s", nodes[j].Pkg, nodes[j].Typ, nodes[j].Method)
		return n1 <= n2
	})
	handleErr(json.NewEncoder(&buf).Encode(nodes))
	handleErr(ioutil.WriteFile("nodes.json", buf.Bytes(), 0o600))

	var graphLoader gfgo.NodeLoader
	err = graphLoader.Load(
		path.Join("github.com/alkemics/goflow/example", "graphs", "..."),
	)
	handleErr(err)

	buf.Reset()

	graphNodes := graphLoader.All()
	pgFilename := "graphs/playground.go"
	goWriter := gfgo.NewWriter(&buf, pgFilename)
	w := gfgo.DebugWriter{
		Writer: goWriter,
	}
	err = gfgo.GeneratePlayground(w, path.Join("github.com/alkemics/goflow/example", "graphs"), graphNodes)
	handleErr(err)

	err = goWriter.Flush()
	handleErr(err)

	err = ioutil.WriteFile(pgFilename, buf.Bytes(), 0o600)
	handleErr(err)
}

func handleErr(errs ...error) {
	nonNilErrs := make([]error, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			nonNilErrs = append(nonNilErrs, err)
		}
	}

	if len(nonNilErrs) == 0 {
		return
	}

	for _, err := range nonNilErrs {
		fmt.Printf("%+v\n", err)
	}
	fmt.Printf("\n%d errors\n", len(errs))
	os.Exit(1)
}

func getWrappers(nodeLoader *gfgo.NodeLoader) []goflow.GraphWrapper {
	return []goflow.GraphWrapper{
		// Linters
		goflow.FromLinter(ids.Lint),

		// Renderers
		inputs.Wrapper,
		goflow.FromNodeWrapper(after.Wrapper),
		gonodes.Wrapper(nodeLoader),
		gonodes.DepWrapper,
		ctx.Wrapper,
		bind.Wrapper,
		goflow.FromNodeWrapper(ifs.Wrapper),
		constants.Wrapper(
			path.Join("github.com/alkemics/goflow", "constants", "..."),
			// Activate the line below to showcase errors when loading constants
			// path.Join("github.com/alkemics/goflow", "notfound", "."),
		),
		outputs.Wrapper,
		varnames.Wrapper,
		types.Wrapper,
		inputs.TypeWrapper,
		outputs.NameWrapper,
		varnames.CompilableWrapper,
		imports.Merger,

		// Checkers
		goflow.FromChecker(cycles.Check),
		goflow.FromChecker(unused.Check),
		goflow.FromChecker(previous.Check),
	}
}
