package constants

import (
	"strconv"
	"strings"

	"github.com/alkemics/goflow"
)

type constant struct {
	name string
	typ  string
	imp  goflow.Import
}

func findConstant(inputName string, constants []constant) constant {
	var cst constant
	for _, c := range constants {
		if c.name == inputName {
			cst = c
		}
	}
	return cst
}

func findHardcodedValue(inputName string) constant {
	var typ string
	if strings.HasPrefix(inputName, "\"") && strings.HasSuffix(inputName, "\"") {
		typ = "string"
	} else if _, err := strconv.ParseFloat(inputName, 64); err == nil {
		// Handles ints and floats
		typ = "@number"
	} else if _, err := strconv.ParseBool(inputName); err == nil {
		// Booleans need to be handled after the numbers because `1` is a valid
		// boolean.
		typ = "bool"
	}
	if typ == "" {
		return constant{}
	}
	return constant{
		name: inputName,
		typ:  typ,
		imp:  goflow.Import{},
	}
}
