package gfgo

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"go/types"
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
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedDeps | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		Tests: true, // Include test files - creates separate test packages with name ending in ".test"
	}, nodesPackages...)
	if err != nil {
		return nil, nil, err
	}

	// Group packages by their base path to merge regular and test packages
	pkgGroups := make(map[string][]*packages.Package)
	for _, pkg := range pkgs {
		// Use PkgPath as the key, stripping any ".test" suffix for test packages
		basePath := strings.TrimSuffix(pkg.PkgPath, ".test")
		pkgGroups[basePath] = append(pkgGroups[basePath], pkg)
	}

	errs := make([]error, 0)
	nodes := make([]Node, 0)
	pkgMap := make(map[string]string)

	for _, group := range pkgGroups {
		// Find the main package (non-test) in the group
		var mainPkg *packages.Package
		for _, pkg := range group {
			if pkg.TypesInfo != nil && !strings.HasSuffix(pkg.PkgPath, ".test") {
				mainPkg = pkg
				break
			}
		}

		if mainPkg == nil {
			continue
		}

		pkgMap[mainPkg.Name] = mainPkg.PkgPath

		// Parse type docs from main package and all test packages, then merge
		typesDoc := make(map[string]string)
		for _, pkg := range group {
			if pkg.Syntax == nil {
				continue
			}
			pkgDoc, err := parseTypesDocFromSyntax(pkg)
			if err != nil {
				errs = append(errs, PkgError{
					PkgPath: mainPkg.PkgPath,
					Err:     err,
				})
				continue
			}
			// Merge docs, preferring non-empty values
			for typeName, doc := range pkgDoc {
				if doc != "" || typesDoc[typeName] == "" {
					typesDoc[typeName] = doc
				}
			}
		}

		// Get constructors first.
		// Start by registering all types.
		allTypes := make(map[string]struct{})
		ignoredTypes := make(map[string]struct{})
		for k, v := range mainPkg.TypesInfo.Defs {
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
		for k, v := range mainPkg.TypesInfo.Defs {
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
					PkgPath: mainPkg.PkgPath,
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

		for k, v := range mainPkg.TypesInfo.Defs {
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
				node, err := createGoNodeFromFunc(decl, signature, mainPkg)
				if err != nil {
					errs = append(errs, PkgError{
						PkgPath: mainPkg.PkgPath,
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
				ns, err := createGoNodeFromType(decl, constructor, named, mainPkg)
				if err != nil {
					errs = append(errs, PkgError{
						PkgPath: mainPkg.PkgPath,
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

// parseTypesDocFromSyntax extracts type documentation from a single package's pre-parsed Syntax.
// This uses pkg.Syntax which is Go 1.27-compatible and avoids re-parsing.
func parseTypesDocFromSyntax(pkg *packages.Package) (map[string]string, error) {
	if pkg.Syntax == nil || len(pkg.Syntax) == 0 {
		return make(map[string]string), nil
	}

	// Create an ast.Package from the pre-parsed Syntax
	astPkg := &ast.Package{
		Name:  pkg.Name,
		Files: make(map[string]*ast.File),
	}

	for _, file := range pkg.Syntax {
		if file != nil && pkg.Fset != nil {
			filename := pkg.Fset.File(file.Pos()).Name()
			astPkg.Files[filename] = file
		}
	}

	if len(astPkg.Files) == 0 {
		return make(map[string]string), nil
	}

	allDecls := doc.New(astPkg, "", doc.AllDecls)
	mappedDoc := make(map[string]string)
	for _, t := range allDecls.Types {
		mappedDoc[t.Name] = t.Doc
	}
	return mappedDoc, nil
}

// parseTypesDoc extracts type documentation from the package's syntax tree.
// Migrated from deprecated parser.ParseDir to use parser.ParseFile (not deprecated) for each file.
func parseTypesDoc(pkg *packages.Package) (map[string]string, error) {
	if len(pkg.GoFiles) == 0 {
		return nil, nil
	}

	// Create an ast.Package by parsing all files (including test files) using parser.ParseFile
	astPkg := &ast.Package{
		Name:  pkg.Name,
		Files: make(map[string]*ast.File),
	}

	// Parse all Go files (pkg.GoFiles includes both regular and test files when Tests: true in loadNodes)
	for _, filename := range pkg.GoFiles {
		file, err := parser.ParseFile(pkg.Fset, filename, nil, parser.ParseComments)
		if err != nil {
			// Skip files that can't be parsed (Go 1.27 compatibility)
			continue
		}
		astPkg.Files[filename] = file
	}

	if len(astPkg.Files) == 0 {
		return nil, PkgError{
			PkgPath: pkg.PkgPath,
			Err:     fmt.Errorf("no valid files could be parsed"),
		}
	}

	allDecls := doc.New(astPkg, "", doc.AllDecls)
	mappedDoc := make(map[string]string)
	for _, t := range allDecls.Types {
		mappedDoc[t.Name] = t.Doc
	}
	return mappedDoc, nil
}
