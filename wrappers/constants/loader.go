package constants

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/packages"

	"github.com/alkemics/goflow"
	"github.com/alkemics/goflow/gfutil/gfgo"
	"github.com/alkemics/lib-go/v9/errors"
)

func loadConstants(constantPackages []string) ([]constant, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo,
	}, constantPackages...)
	if err != nil {
		return nil, err
	}

	errs := make([]error, 0)
	constants := make([]constant, 0)
	for _, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			errs = append(errs, craftPkgError(pkg.PkgPath, pkg.Errors))
		}

		if pkg.TypesInfo == nil {
			continue
		}

		for k, v := range pkg.TypesInfo.Defs {
			if v == nil || k.Obj == nil {
				continue
			}

			// We are only interested in exported variables/constants.
			if !ast.IsExported(k.Name) || k.Obj.Kind != ast.Con && k.Obj.Kind != ast.Var {
				continue
			}

			typ, _, err := gfgo.ParseType(v.Type())
			if err != nil {
				return nil, errors.New(
					"could not find type of constant {{constant}}",
					errors.WithField("constant", k.Name),
					errors.WithCause(err),
				)
			}
			constants = append(constants, constant{
				name: fmt.Sprintf("%s.%s", pkg.Name, k.Name),
				typ:  typ,
				imp: goflow.Import{
					Pkg: pkg.Name,
					Dir: pkg.PkgPath,
				},
			})
		}
	}

	if len(errs) > 0 {
		err = goflow.MultiError{Errs: errs}
	}

	return constants, err
}
