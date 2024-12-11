package main

import (
	"encoding/json"
	"fmt"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

type Node struct {
	Name string
	Type string
}

type Edge struct {
	From string
	To   string
}

type Graph struct {
	Nodes []Node
	Edges []Edge
}

func Reverse[T any](ss []T) []T {
	if len(ss) < 2 {
		return ss
	}

	sorted := make([]T, len(ss))
	for i := 0; i < len(ss); i++ {
		sorted[i] = ss[len(ss)-i-1]
	}

	return sorted
}

func main() {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		fmt.Println("error reading go.mod:", err)
		return
	}

	modFile, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		fmt.Println("Error parsing go.mod:", err)
		return
	}

	packages := make(map[string][]Node)
	packagesStack := make([]string, 0)
	visited := make(map[string]bool)
	analyzePackage(modFile.Module.Mod.Path, visited, packages, &packagesStack)

	nodes := make([]Node, 0)
	edges := make([]Edge, 0)
	for _, v := range Reverse(packagesStack) {
		for _, vv := range packages[v] {
			isExistNode := false
			for _, n := range nodes {
				if n == vv {
					isExistNode = true
				}
			}
			if !isExistNode {
				nodes = append(nodes, vv)
			}

			edges = append(edges, Edge{From: v, To: vv.Name})
		}
		nodes = append(nodes, Node{Name: v, Type: "loc"})
	}

	gg := Graph{Nodes: nodes, Edges: edges}
	// fmt.Println(packagesStack)
	// fmt.Println(packages)

	b, _ := json.Marshal(gg)
	fmt.Println(string(b))
}

func analyzePackage(importPath string, visited map[string]bool, packages map[string][]Node, packagesStack *[]string) {
	if visited[importPath] {
		return
	}
	visited[importPath] = true

	pkg, err := build.Import(importPath, "", 0)
	if err != nil {
		log.Printf("error importing package %s: %v", importPath, err)
		return
	}

	for _, imp := range pkg.Imports {
		pkg, err := build.Import(imp, "", 0)
		if err != nil {
			log.Printf("error importing package %s: %v", importPath, err)
			return
		}

		packageType := "ext"
		if pkg.Goroot {
			packageType = "std"
		}
		if isLocaclPkg(pkg.Dir) {
			packageType = "loc"
		}

		packages[importPath] = append(packages[importPath], Node{Name: imp, Type: packageType})

		isExistPkg := false
		for _, v := range *packagesStack {
			if v == importPath {
				isExistPkg = true
			}
		}
		if !isExistPkg {
			*packagesStack = append(*packagesStack, importPath)
		}

		if packageType != "loc" {
			continue
		}

		analyzePackage(imp, visited, packages, packagesStack)
	}
}

func isLocaclPkg(pkgDir string) bool {
	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current directory: %v", err)
	}

	absPkgDir, err := filepath.Abs(pkgDir)
	if err != nil {
		log.Fatalf("failed to get absolute path of package: %v", err)
	}

	return strings.HasPrefix(absPkgDir, projectRoot)
}
