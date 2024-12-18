package graph

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/elliotchance/pie/v2"
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

func New() Graph {
	packages := make(map[string][]Node)
	timelinePackages := make([]string, 0)

	modPath := getModPath()
	analyzePackage(modPath, packages, &timelinePackages)
	nodes, edges := getNodesAndEdges(packages, timelinePackages)

	return Graph{Nodes: nodes, Edges: edges}
}

func getModPath() string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		fmt.Println("error reading go.mod:", err)
		return ""
	}

	modFile, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		fmt.Println("error parsing go.mod:", err)
		return ""
	}

	return modFile.Module.Mod.Path
}

func analyzePackage(importPath string, packages map[string][]Node, timelinePackages *[]string) {
	if _, ok := packages[importPath]; ok {
		return
	}

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
		if isLocalPkg(pkg.Dir) {
			packageType = "loc"
		}

		packages[importPath] = append(packages[importPath], Node{Name: imp, Type: packageType})

		isExistPkg := false
		for _, v := range *timelinePackages {
			if v == importPath {
				isExistPkg = true
			}
		}
		if !isExistPkg {
			*timelinePackages = append(*timelinePackages, importPath)
		}

		if packageType != "loc" {
			continue
		}

		analyzePackage(imp, packages, timelinePackages)
	}
}

func isLocalPkg(pkgDir string) bool {
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

func getNodesAndEdges(packages map[string][]Node, timelinePackages []string) ([]Node, []Edge) {
	nodes := make([]Node, 0)
	edges := make([]Edge, 0)
	for _, v := range pie.Reverse(timelinePackages) {
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

	return nodes, edges
}
