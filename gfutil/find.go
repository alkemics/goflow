package gfutil

import (
	"io/ioutil"
	"path"
	"strings"
)

func FindGraphFileNames(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	graphNames := make([]string, 0)

	for _, file := range files {
		if file.IsDir() && file.Name() != "vendor" {
			subs, err := FindGraphFileNames(path.Join(dir, file.Name()))
			if err != nil {
				return nil, err
			}

			graphNames = append(graphNames, subs...)
		} else if strings.HasSuffix(file.Name(), ".yml") {
			graphNames = append(graphNames, path.Join(dir, file.Name()))
		}
	}

	return graphNames, nil
}
