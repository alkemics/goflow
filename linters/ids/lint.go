package ids

import (
	"fmt"
	"strings"
)

type Error struct {
	DuplicatedIDs []string
}

func (e Error) Error() string {
	return fmt.Sprintf("duplicated not ids: [%s] (duplicated-ids)", strings.Join(e.DuplicatedIDs, ", "))
}

func Lint(unmarshal func(interface{}) error) error {
	var graph struct {
		Nodes []struct {
			ID string `yaml:"id"`
		} `yaml:"nodes"`
	}
	if err := unmarshal(&graph); err != nil {
		return err
	}

	nodeIDMap := make(map[string]struct{})
	duplicatedIDMap := make(map[string]struct{})
	for _, n := range graph.Nodes {
		if _, ok := nodeIDMap[n.ID]; ok {
			duplicatedIDMap[n.ID] = struct{}{}
		}
		nodeIDMap[n.ID] = struct{}{}
	}

	if len(duplicatedIDMap) > 0 {
		duplicatedIDs := make([]string, 0)
		for id := range duplicatedIDMap {
			duplicatedIDs = append(duplicatedIDs, id)
		}
		return Error{
			DuplicatedIDs: duplicatedIDs,
		}
	}

	return nil
}
