package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/stepan2volkov/godepuml/godepuml"
	"golang.org/x/mod/modfile"
)

var (
	outputPath = flag.String("o", "diagram.puml", "Output file")
	path       = flag.String("p", "./", "Root of your Go project")
)

func main() {
	flag.Parse()

	moduleName := getModuleName(*path)

	absPath, err := filepath.Abs(*path)
	if err != nil {
		fmt.Println("path is invalid:", err)
		os.Exit(1)
	}

	scanner := godepuml.PackageScanner{
		Root:       absPath,
		ModuleName: moduleName,
	}

	pkgList, err := scanner.Scan(absPath)
	if err != nil {
		log.Fatalln(err)
	}

	f, err := os.Create(*outputPath)
	if err != nil {
		fmt.Printf("error while creating file '%s': %v", *outputPath, err)
		os.Exit(1)
	}

	fmt.Fprintf(f, "@startuml '%s'\n\n", moduleName)

	for _, pkg := range pkgList {
		for dep := range pkg.Deps {
			fmt.Fprintf(f, "[%s] --> [%s]\n", pkg.Path, dep)
		}
	}

	fmt.Fprintln(f)
	fmt.Fprintln(f, "@enduml")
}

func getModuleName(path string) string {
	content, err := os.ReadFile(filepath.Join(path, "go.mod"))
	if err != nil {
		fmt.Printf("error when getting go.mod: %v\n", err)
		os.Exit(1)
	}

	modName := modfile.ModulePath(content)

	return modName
}
