package gfgo

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/alkemics/goflow"
	"github.com/alkemics/lib-go/v9/errors"
)

func ParseType(typ types.Type) (string, []goflow.Import, error) {
	imports, err := extractImports(typ)
	if err != nil {
		return "", nil, err
	}
	// Replace the full packages names with their name.
	t := typ.String()
	for i, imp := range imports {
		t = strings.Replace(t, imp.Dir, imp.Pkg, -1)
		// TODO: check this when we migrate to modules.
		if strings.Contains(imp.Dir, "/vendor/") {
			imports[i].Dir = imp.Dir[strings.Index(imp.Dir, "/vendor/")+8:]
		}
	}
	return t, imports, nil
}

func extractImports(typ types.Type) ([]goflow.Import, error) {
	switch typ := typ.(type) {
	case *types.Basic:
		return nil, nil
	case *types.Interface:
		return nil, nil
	case *types.Array:
		return extractImports(typ.Elem())
	case *types.Slice:
		return extractImports(typ.Elem())
	case *types.Pointer:
		return extractImports(typ.Elem())
	case *types.Named:
		pkg := typ.Obj().Pkg()
		if pkg == nil {
			// Typically the case for built-ins like 'error'
			return nil, nil
		}

		return []goflow.Import{
			{
				Pkg: pkg.Name(),
				Dir: pkg.Path(),
			},
		}, nil
	case *types.Struct:
		imports, _, err := toFields(structFields(typ))
		if err != nil {
			return nil, err
		}

		return imports, nil
	case *types.Signature:
		imports, _, _, err := ParseSignature(typ)
		if err != nil {
			return nil, err
		}

		return imports, nil
	case *types.Map:
		keyImports, err := extractImports(typ.Key())
		if err != nil {
			return nil, err
		}

		elemImports, err := extractImports(typ.Elem())
		if err != nil {
			return nil, err
		}

		imports := append(keyImports, elemImports...)
		return imports, nil
	}
	return nil, errors.New("could no find type of {{type}}", errors.WithField("type", fmt.Sprintf("%T", typ)))
}

func ParseSignature(signature *types.Signature) (imports []goflow.Import, inputs, outputs []goflow.Field, err error) {
	var imps []goflow.Import

	// Parse inputs.
	if signature.Params() != nil {
		imps, inputs, err = toFields(signature.Params())
		if err != nil {
			return nil, nil, nil, errors.New("error reading inputs", errors.WithCause(err))
		}

		imports = append(imports, imps...)
	}

	// Parse outputs.
	if signature.Results() != nil {
		imps, outputs, err = toFields(signature.Results())
		if err != nil {
			return nil, nil, nil, errors.New("error reading outputs", errors.WithCause(err))
		}

		imports = append(imports, imps...)
	}

	return imports, inputs, outputs, nil
}

func toFields(vars *types.Tuple) ([]goflow.Import, []goflow.Field, error) {
	fields := make([]goflow.Field, vars.Len())
	mergedImports := make([]goflow.Import, 0)
	for i := 0; i < vars.Len(); i++ {
		v := vars.At(i)
		typ, imports, err := ParseType(v.Type())
		if err != nil {
			return nil, nil, errors.New(
				"could not read type of input {{input}}",
				errors.WithField("input", fmt.Sprintf("#%d", i)),
				errors.WithCause(err),
			)
		}

		fields[i] = goflow.Field{
			Name: v.Name(),
			Type: typ,
		}

		mergedImports = append(mergedImports, imports...)
	}
	return mergedImports, fields, nil
}

func structFields(s *types.Struct) *types.Tuple {
	vars := make([]*types.Var, s.NumFields())
	for i := 0; i < s.NumFields(); i++ {
		vars[i] = s.Field(i)
	}
	return types.NewTuple(vars...)
}
