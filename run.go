package goflow

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

// WriteFile writes content into a file called filename. It updates the file only
// if there is a difference between the current content of the file and what is
// passed in content.
//
// It also logs the time taken since start, because why not.
func WriteFile(content, filename string, start time.Time) error {
	before, err := ioutil.ReadFile(filename)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if string(before) == content {
		// Do not write file if there is nothing new to write...
		return nil
	}

	created := os.IsNotExist(err)

	wf, err := os.OpenFile(filename, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		return err
	}
	defer wf.Close()

	if _, err := wf.WriteString(content); err != nil {
		return err
	}

	action := "updated"
	if created {
		action = "created"
	}
	fmt.Println(filename, action, "in", time.Since(start))

	return nil
}

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
