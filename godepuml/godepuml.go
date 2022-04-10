package godepuml

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type Package struct {
	Path string
	Deps map[string]struct{}
}

type PackageList map[string]Package

func (l1 PackageList) append(p Package) {
	path := p.Path
	if _, found := l1[path]; !found {
		l1[path] = p
		return
	}
	for dep := range p.Deps {
		l1[path].Deps[dep] = struct{}{}
	}
}

func (l1 PackageList) merge(l2 PackageList) {
	for _, pkg := range l2 {
		l1.append(pkg)
	}
}

type PackageScanner struct {
	Root         string
	ModuleName   string
	ExcludedDirs map[string]struct{}
}

func (s *PackageScanner) Scan(path string) (PackageList, error) {
	ret := make(PackageList)

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf(`error when open dir "%s" (%w)`, path, err)
	}

	pkgPath := strings.TrimPrefix(path, s.Root)
	pkgPath = strings.TrimPrefix(pkgPath, "/")
	pkgPath = strings.ReplaceAll(pkgPath, "/", ".")

	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())

		if _, found := s.ExcludedDirs[entryPath]; found {
			continue
		}

		if entry.IsDir() {
			pkgs, err := s.Scan(entryPath)
			if err != nil {
				return nil, err
			}
			ret.merge(pkgs)
		}

		pkg, err := s.getPkg(pkgPath, entryPath)
		if err != nil {
			return nil, err
		}

		ret.append(pkg)
	}

	return ret, nil
}

func (s *PackageScanner) getPkg(pkgPath, filepath string) (Package, error) {
	if strings.HasSuffix(filepath, "_test.go") {
		return Package{}, nil
	}
	if !strings.HasSuffix(filepath, ".go") {
		return Package{}, nil
	}

	fset := token.NewFileSet()

	fi, err := parser.ParseFile(fset, filepath, nil, parser.ImportsOnly)
	if err != nil {
		return Package{}, fmt.Errorf(`error when getting imports for %s: %w`, filepath, err)
	}

	deps := make(map[string]struct{}, len(fi.Imports))

	for _, importSpec := range fi.Imports {
		dep := strings.Trim(importSpec.Path.Value, `"`)

		if strings.HasPrefix(dep, s.ModuleName) {
			dep = strings.TrimPrefix(dep, s.ModuleName)
			dep = strings.TrimPrefix(dep, "/")
			dep = strings.ReplaceAll(dep, "/", ".")
			deps[dep] = struct{}{}
		}
	}

	return Package{
		Path: pkgPath,
		Deps: deps,
	}, nil
}
