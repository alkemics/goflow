package gfgo

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"go/types"
	"path"
	"strings"
	"sync"

	"golang.org/x/tools/go/packages"

	"github.com/alkemics/goflow"
)

type Node struct {
	Pkg          string
	Typ          string
	PkgPath      string
	Doc          string
	Constructor  string
	Method       string
	Filename     string
	Imports      []goflow.Import
	Dependencies []goflow.Field
	Inputs       []goflow.Field
	Outputs      []goflow.Field
}

func (n Node) Match(typ string) bool {
	if n.Method != "" && n.Method != "Run" {
		return typ == fmt.Sprintf("%s.%s.%s", n.Pkg, n.Typ, n.Method)
	}
	return typ == fmt.Sprintf("%s.%s", n.Pkg, n.Typ)
}

type NodeLoader struct {
	mu    sync.Mutex
	pkgs  map[string]string
	nodes []Node
}

func (l *NodeLoader) Load(pkgs ...string) error {
	nodes, pkgMap, err := loadNodes(pkgs)
	if err != nil {
		return err
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.pkgs == nil {
		l.pkgs = make(map[string]string)
	}

	for k, v := range pkgMap {
		l.pkgs[k] = v
	}

	l.nodes = nodes

	return nil
}

func (l *NodeLoader) Refresh(pkgName string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	pkgPath := l.pkgs[pkgName]
	if pkgPath == "" {
		return fmt.Errorf("unkown pkg: %s", pkgName)
	}

	nodes, _, err := loadNodes([]string{pkgPath})
	if err != nil {
		return err
	}

	otherNodes := make([]Node, 0, len(l.nodes))

	for _, node := range l.nodes {
		if node.PkgPath != pkgPath {
			otherNodes = append(otherNodes, node)
		}
	}

	l.nodes = append(otherNodes, nodes...)

	return nil
}

func (l *NodeLoader) All() []Node {
	l.mu.Lock()
	defer l.mu.Unlock()

	nodes := make([]Node, len(l.nodes))
	copy(nodes, l.nodes)

	return nodes
}

func (l *NodeLoader) Find(typ string) Node {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, node := range l.nodes {
		if node.Match(typ) {
			return node
		}
	}

	return Node{}
}

type typeConstructor struct {
	name         string
	imports      []goflow.Import
	dependencies []goflow.Field
}

func loadNodes(nodesPackages []string) ([]Node, map[string]string, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedDeps | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
	}, nodesPackages...)
	if err != nil {
		return nil, nil, err
	}

	errs := make([]error, 0)
	nodes := make([]Node, 0)
	pkgMap := make(map[string]string)
	for _, pkg := range pkgs {
		if pkg.TypesInfo == nil {
			continue
		}

		pkgMap[pkg.Name] = pkg.PkgPath

		typesDoc, err := parseTypesDoc(pkg)
		if err != nil {
			errs = append(errs, PkgError{
				PkgPath: pkg.PkgPath,
				Err:     err,
			})
		}

		// Get constructors first.
		// Start by registering all types.
		allTypes := make(map[string]struct{})
		ignoredTypes := make(map[string]struct{})
		for k, v := range pkg.TypesInfo.Defs {
			if v == nil ||
				!v.Exported() ||
				k.Obj == nil ||
				k.Obj.Decl == nil ||
				k.Obj.Kind != ast.Typ {
				continue
			}
			decl, ok := k.Obj.Decl.(*ast.TypeSpec)
			if !ok {
				continue
			}
			typ := decl.Name.String()
			allTypes[typ] = struct{}{}
			if shouldIgnoreNode(typesDoc[typ]) {
				ignoredTypes[typ] = struct{}{}
			}
		}

		// Then register all constructors by type.
		constructorNames := make(map[string]struct{})
		constructors := make(map[string]typeConstructor)
		for k, v := range pkg.TypesInfo.Defs {
			if v == nil ||
				!v.Exported() ||
				k.Obj == nil ||
				k.Obj.Decl == nil ||
				k.Obj.Kind != ast.Fun ||
				!strings.HasPrefix(k.Obj.Name, "New") {
				continue
			}

			signature, ok := v.Type().(*types.Signature)
			if !ok || signature.Results() == nil {
				continue
			}

			decl, ok := k.Obj.Decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			imports, dependencies, results, err := ParseSignature(signature)
			if err != nil {
				errs = append(errs, PkgError{
					PkgPath: pkg.PkgPath,
					Err:     fmt.Errorf("parsing %s: %w", decl.Name.String(), err),
				})
			}

			if len(results) != 1 {
				continue
			}

			// Get the type of returned by the constructor.
			typ := strings.TrimPrefix(results[0].Type, "*")
			split := strings.Split(typ, ".")
			if len(split) > 1 {
				typ = split[len(split)-1]
			}

			if _, ok := allTypes[typ]; !ok {
				// It's not a constructor.
				continue
			}

			constructorNames[decl.Name.String()] = struct{}{}
			if _, ok := ignoredTypes[typ]; ok {
				continue
			}

			constructors[typ] = typeConstructor{
				name:         decl.Name.String(),
				imports:      imports,
				dependencies: dependencies,
			}
		}

		for k, v := range pkg.TypesInfo.Defs {
			if v == nil || !v.Exported() || k.Obj == nil || k.Obj.Decl == nil {
				continue
			}
			switch k.Obj.Kind {
			case ast.Fun:
				// Load the functions (we don't have the methods here).
				signature, ok := v.Type().(*types.Signature)
				if !ok {
					continue
				}
				decl, ok := k.Obj.Decl.(*ast.FuncDecl)
				if !ok {
					continue
				}
				if _, ok := constructorNames[decl.Name.String()]; ok {
					// Ignore constructors.
					continue
				}
				if shouldIgnoreNode(decl.Doc.Text()) {
					continue
				}
				node, err := createGoNodeFromFunc(decl, signature, pkg)
				if err != nil {
					errs = append(errs, PkgError{
						PkgPath: pkg.PkgPath,
						Err:     fmt.Errorf("parsing %s: %w", decl.Name, err),
					})
				} else {
					nodes = append(nodes, node)
				}
			case ast.Typ:
				// Load the types with their methods.
				named, ok := v.Type().(*types.Named)
				if !ok {
					continue
				}
				decl, ok := k.Obj.Decl.(*ast.TypeSpec)
				if !ok {
					continue
				}
				typ := decl.Name.String()
				if _, ok := ignoredTypes[typ]; ok {
					continue
				}
				constructor := constructors[typ]
				if constructor.name == "" {
					continue
				}
				ns, err := createGoNodeFromType(decl, constructor, named, pkg)
				if err != nil {
					errs = append(errs, PkgError{
						PkgPath: pkg.PkgPath,
						Err:     fmt.Errorf("parsing %s: %w", decl.Name, err),
					})
				} else {
					nodes = append(nodes, ns...)
				}
			default:
				continue
			}
		}
	}

	if len(errs) > 0 {
		err = goflow.MultiError{Errs: errs}
	}

	return nodes, pkgMap, err
}

func createGoNodeFromFunc(decl *ast.FuncDecl, signature *types.Signature, pkg *packages.Package) (Node, error) {
	imports, inputs, outputs, err := ParseSignature(signature)
	if err != nil {
		return Node{}, err
	}

	return Node{
		Pkg:     pkg.Name,
		Typ:     decl.Name.String(),
		PkgPath: pkg.PkgPath,
		Doc:     decl.Doc.Text(),
		// Add node import.
		Imports: append(
			imports,
			goflow.Import{
				Pkg: pkg.Name,
				Dir: pkg.PkgPath,
			},
		),
		Inputs:   inputs,
		Outputs:  outputs,
		Filename: pkg.Fset.File(decl.Pos()).Name(),
	}, nil
}

func createGoNodeFromType(decl *ast.TypeSpec, constructor typeConstructor, named *types.Named, pkg *packages.Package) ([]Node, error) {
	baseImports := append(constructor.imports, goflow.Import{
		Pkg: pkg.Name,
		Dir: pkg.PkgPath,
	})

	// Add one node per exported method.
	errs := make([]error, 0)
	nodes := make([]Node, 0, named.NumMethods())
	for i := 0; i < named.NumMethods(); i++ {
		method := named.Method(i)
		if !method.Exported() {
			continue
		}

		doc := getDocFromPackage(pkg, method.Pos())
		if shouldIgnoreNode(doc) {
			continue
		}

		signature, ok := method.Type().(*types.Signature)
		if !ok {
			continue
		}

		imports, inputs, outputs, err := ParseSignature(signature)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		nodes = append(nodes, Node{
			Pkg:     pkg.Name,
			Typ:     decl.Name.String(),
			PkgPath: pkg.PkgPath,
			Doc:     doc,
			// TODO: parse deps from struct rather than from constructor if not supplied
			Dependencies: constructor.dependencies,
			Constructor:  constructor.name,
			Method:       method.Name(),
			Imports:      append(imports, baseImports...),
			Inputs:       inputs,
			Outputs:      outputs,
			Filename:     pkg.Fset.File(decl.Pos()).Name(),
		})
	}

	var err error
	if len(errs) > 0 {
		err = goflow.MultiError{Errs: errs}
	}

	return nodes, err
}

func shouldIgnoreNode(doc string) bool {
	for _, line := range strings.Split(doc, "\n") {
		if line == "node:ignore" {
			return true
		}
	}
	return false
}

// getDocFromPackage loads the doc directly from the syntax if we don't have it already.
func getDocFromPackage(pkg *packages.Package, pos token.Pos) string {
	for _, f := range pkg.Syntax {
		for _, d := range f.Decls {
			switch decl := d.(type) {
			case *ast.FuncDecl:
				if decl.Name.NamePos == pos {
					return decl.Doc.Text()
				}
			}
		}
	}
	return ""
}

// parseTypesDoc hacks the thing by grabbing info of types using doc.
// TODO: review all this and improve if possible...
//
//	shall we try with https://pkg.go.dev/go/ast?tab=doc#CommentGroup maybe?
func parseTypesDoc(pkg *packages.Package) (map[string]string, error) {
	if len(pkg.GoFiles) == 0 {
		return nil, nil
	}

	pkgs, err := parser.ParseDir(pkg.Fset, path.Dir(pkg.GoFiles[0]), nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	if _, ok := pkgs[pkg.Name]; !ok {
		return nil, PkgError{
			PkgPath: pkg.PkgPath,
			Err:     fmt.Errorf("package %s not found", pkg.Name),
		}
	}

	allDecls := doc.New(pkgs[pkg.Name], "", doc.AllDecls)
	mappedDoc := make(map[string]string)
	for _, t := range allDecls.Types {
		mappedDoc[t.Name] = t.Doc
	}
	return mappedDoc, nil
}
