package gfgo

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/alkemics/goflow"
)

type jsonNode struct {
	ID           string         `json:"id"`
	Pkg          string         `json:"pkg"`
	Type         string         `json:"type"`
	Inputs       []goflow.Field `json:"inputs"`
	Outputs      []goflow.Field `json:"outputs"`
	Dependencies []goflow.Field `json:"dependencies"`
}

type jsonEdge struct {
	SourceID string `json:"sourceId"`
	TargetID string `json:"targetId"`
	Type     string `json:"inputType"` // json tag to have it work with the current dashboard
}

type bindings []string

func (b *bindings) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multipleBindings []string
	if err := unmarshal(&multipleBindings); err == nil {
		*b = multipleBindings
		return nil
	}

	var simpleBinding string
	if err := unmarshal(&simpleBinding); err != nil {
		return err
	}

	*b = []string{simpleBinding}
	return nil
}

// WithJSONMarshal wraps the generate function to add a MarshalJSON function
// to the graph.
//
// TODO:
// This method does some work similar to the bind, if and after wrappers, there
// probably is a simple way to handle that without actually redoing the work.
// I tried coalescing the nodes from the bind wrapper, it did not go too well :p
func WithJSONMarshal(
	generate GenerateFunc,
	yamlFilename string,
	nodes []Node,
) GenerateFunc {
	file, err := os.Open(yamlFilename)
	if err != nil {
		return func(w io.Writer, g goflow.GraphRenderer) error {
			return err
		}
	}

	var graph struct {
		Name  string
		Nodes []struct {
			ID    string
			Type  string
			Bind  map[string]bindings
			If    []string
			After []string
		}
		Outputs map[string]bindings
	}
	if err := yaml.NewDecoder(file).Decode(&graph); err != nil {
		return func(w io.Writer, g goflow.GraphRenderer) error {
			return goflow.ParseYAMLError(err)
		}
	}

	jsonNodes := make([]jsonNode, len(graph.Nodes))
	edges := make([]jsonEdge, 0)
	for i, n := range graph.Nodes {
		jn := jsonNode{
			ID: n.ID,
		}

		for _, node := range nodes {
			if node.Match(n.Type) {
				jn.Type = node.Typ
				jn.Pkg = node.Pkg
				jn.Inputs = node.Inputs
				jn.Outputs = node.Outputs
				jn.Dependencies = node.Dependencies
				break
			}
		}

		jsonNodes[i] = jn

		for inputName, bds := range n.Bind {
			var typ string
			for _, input := range jn.Inputs {
				if inputName == input.Name {
					typ = input.Type
				}
			}

			for _, b := range bds {
				var sourceID string
				if split := strings.SplitN(b, ".", 2); len(split) == 2 {
					sID := split[0]
					if sID == "inputs" {
						sourceID = sID
					} else {
						for _, on := range graph.Nodes {
							if on.ID == sID {
								sourceID = sID
							}
						}
					}
				}

				if sourceID == "" {
					sourceID = "manualInput"
				}

				edges = append(edges, jsonEdge{
					SourceID: sourceID,
					TargetID: n.ID,
					Type:     typ,
				})
			}
		}

		for _, i := range n.If {
			i = strings.TrimPrefix(i, "not")
			i = strings.TrimSpace(i)
			var sourceID string
			if split := strings.SplitN(i, ".", 2); len(split) == 2 {
				sID := split[0]
				if sID == "inputs" {
					sourceID = sID
				} else {
					for _, on := range graph.Nodes {
						if on.ID == sID {
							sourceID = sID
						}
					}
				}
			}

			if sourceID == "" {
				sourceID = "manualInput"
			}

			edges = append(edges, jsonEdge{
				SourceID: sourceID,
				TargetID: n.ID,
				Type:     "bool",
			})
		}

		for _, a := range n.After {
			edges = append(edges, jsonEdge{
				SourceID: a,
				TargetID: n.ID,
				Type:     "__wait__",
			})
		}
	}

	for _, bds := range graph.Outputs {
		for _, b := range bds {
			var sourceID string
			if split := strings.SplitN(b, ".", 2); len(split) == 2 {
				sID := split[0]
				if sID == "inputs" {
					sourceID = sID
				} else {
					for _, on := range graph.Nodes {
						if on.ID == sID {
							sourceID = sID
						}
					}
				}
			}

			if sourceID == "" {
				sourceID = "manualInput"
			}

			edges = append(edges, jsonEdge{
				SourceID: sourceID,
				TargetID: "outputs",
				Type:     "whatever",
			})
		}
	}

	// Check all the nodes that have not been used
	for _, n := range jsonNodes {
		used := false
		for _, e := range edges {
			if n.ID == e.SourceID {
				used = true
				break
			}
		}

		if !used {
			edges = append(edges, jsonEdge{
				SourceID: n.ID,
				TargetID: "outputs",
				Type:     "__wait__",
			})
		}
	}

	sort.SliceStable(jsonNodes, func(i, j int) bool {
		return jsonNodes[i].ID <= jsonNodes[j].ID
	})
	sort.SliceStable(edges, func(i, j int) bool {
		return fmt.Sprintf("%v", edges[i]) <= fmt.Sprintf("%v", edges[j])
	})

	return func(w io.Writer, graph goflow.GraphRenderer) error {
		if err := generate(w, graph); err != nil {
			return err
		}

		msg, err := json.Marshal(map[string]interface{}{
			"id":   graph.Name(),
			"type": graph.Name(),
			"pkg":  graph.Pkg(),

			"filename": yamlFilename,
			"doc":      graph.Doc(),

			"nodes": jsonNodes,
			"edges": edges,

			"inputs":  graph.Inputs(),
			"outputs": graph.Outputs(),
		})
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w, `
// MarshalJSON returns the json representation of the graphs. It is pre-generated by
// WithJSONMarshal, and hence never returns an error.
func (g %s) MarshalJSON() ([]byte, error) {
	return []byte(%s), nil
}
`,
			graph.Name(),
			strconv.Quote(string(msg)),
		)
		return err
	}
}
